package task

import (
	"context"
	"io"
	"log"
	"os"
	"time"
    "math"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

// A State type represents the states a task goes through.
// From Pending, Scheduled, Running to Failed or Completed.
type State int

const (
    Pending State = iota
    Scheduled
    Running
    Completed
    Failed
)

// Attributes a task should have are UUIDs and Name
// UUID will uniquely identify an individual task 
// Name will serve as a human-readable reference of the task
//
// For this project, we will deal with docker containers
// Image will be the image the container will use
// CPU, Memory and Disk will help the system to identify the ammount of resources a task needs
// ExposedPorts and PortBindings will be used to properly allocate network ports for the task
// and  network ports available on the network.
// RestartPolicy will tell the system what to do in the event of a task stopping or failling unexpectedly.
//
// To know when a task stops, we need StartTime and EndTime
// And to tie to a container after its launch, we add ContainerId
//
// TODO: Move ContainerId, Cpu, Image, Memory, Disk, ExposedPorts and PortBindings to a separate type
type Task struct {
    Id             uuid.UUID
    ContainerId    string
    Name           string
    State          State
    Cpu            float64
    Image          string
    Memory         int
    Disk           int
    ExposedPorts   nat.PortSet 
    PortBindings   map[string]string
    RestartPolicy  string
    StartTime      time.Time
    EndTime        time.Time
}

// To tell the system how to stop a Task, we need to create TaskEvent
// And similarly to Task, TaskEvent will need it's own identifier and State,
// which will indicate the state the task should transition (eg. from Running to Completed)
// along with a Timestamp to record the time the event was requested
//
// A TaskEvent will contain a Task struct and the TaskEvent will be an internal object
// the system will use to trigger tasks from one state to another
type TaskEvent struct {
    Id             uuid.UUID
    State          State
    Timestamp      time.Time
    Task           Task
}

// A Config type will be used to create a Task
// It will contain all the information needed to create a Task
//
// - Name -> will be used to identify the task
// - AttachStdin, AttachStdout and AttachStderr --> will be used to attach the task to the system's standard streams
// - ExposedPorts --> will be used to expose ports to the network
// - Cmd --> will be used to run a command inside the container
// - Image --> will be the image the container will use
// - Cpu, Memory and Disk --> will help the system to identify the ammount of resources a task needs
//                            and also will serve the purpose of finiding a node 
//                            in the cluster that can run the task
// - Env --> will be used to set environment variables inside the container
// - RestartPolicy --> will tell the system what to do in the event of a task stopping or 
//                     failling unexpectedly. This field is one of the mechanisms to ensure 
//                     some resiliency to the system
type Config struct {
    Name string
    AttachStdin bool
    AttachStdout bool
    AttachStderr bool
    ExposedPorts nat.PortSet
    Cmd []string
    Image string
    Cpu float64
    Memory int64
    Disk int64
    Env []string
    RestartPolicy string
}

// Docker struct will encapsulate everything needed to run a Docker task
// The Client field will hold a Docker client object used to interact with the Docker API
// The Config field wll hold all the configuration for the task
type Docker struct {
    Client *client.Client
    Config Config
}

// DockerResult struct is used to return the result of a Docker task
// It will hold the followint fields:
// - Error --> will hold any error that might happens 
// - Action --> will hold the action that was performed (eg. create, start, stop, etc)
// - ContainerId --> will hold the container id of the task
// - Result --> will hold the result of the action (eg. success, failure)
type DockerResult struct {
    Error error
    Action string
    ContainerId string
    Result string
}
    
// The Run method will pull the Docker image the task will use fron a given registry
// To do this, it first creates a context and then uses the Docker client to pull the image
// passing the context, image name and any options needed to pull the image
func (d *Docker) Run() DockerResult {
    ctx := context.Background()
    reader, err := d.Client.ImagePull(
        ctx, 
        d.Config.Image, 
        types.ImagePullOptions{},
    )

    if err != nil {
        log.Printf("Error while pulling image %s: %v\n", d.Config.Image, err)
        return DockerResult{Error: err,}
    }
    // Copy the output from the image pull to stdout
    io.Copy(os.Stdout, reader)

    // After pulling the image, the Run method will create a container
    // To do this, it first creates a context and then uses the Docker client to create the container
    // passing the context, container name and any options needed to create the container
    //
    // Create arguments to create a container
    restart_policy := container.RestartPolicy{
        Name: d.Config.RestartPolicy,
    }
    resources := container.Resources{
        Memory: d.Config.Memory,
        NanoCPUs: int64(d.Config.Cpu * math.Pow(10,9)),
    }
    container_config := container.Config{
        Image: d.Config.Image,
        Tty: false,
        Env: d.Config.Env,
        ExposedPorts: d.Config.ExposedPorts,
    }

    // HostConfig will hold the configuration a task requires of the host on which 
    // it will run and expose automatically the ports needed by the task 
    host_config := container.HostConfig{
        RestartPolicy: restart_policy,
        Resources: resources,
        PublishAllPorts: true,
    }

    // Create the container
    resp, err := d.Client.ContainerCreate(
        ctx,
        &container_config,
        &host_config,
        nil, // NetworkingConfig
        nil, // Platform
        d.Config.Name,
    )
    if err != nil {
        log.Printf("Error while creating container using image %s: %v\n", d.Config.Image, err)
        return DockerResult{Error: err,}
    }

    // Start the container
    err = d.Client.ContainerStart(
        ctx,
        resp.ID,
        types.ContainerStartOptions{},
    )
    if err != nil {
        log.Printf("Error while starting container %s: %v\n", resp.ID, err)
        return DockerResult{Error: err,}
    }

    // Serve the logs of the container
    out, err := d.Client.ContainerLogs(
        ctx,
        resp.ID,
        types.ContainerLogsOptions{
            ShowStdout: true,
            ShowStderr: true,
        },
    )
    if err != nil {
        log.Printf("Error while getting logs from container %s: %v\n", resp.ID, err)
        return DockerResult{Error: err,}
    }

    // Copy the output from the container logs to stdout and stderr
    stdcopy.StdCopy(os.Stdout, os.Stderr, out)
    return DockerResult{
        ContainerId: resp.ID,
        Action: "start",
        Result: "success",
        Error: nil,
    }
}


// The Stop method will stop a Docker container
func (d *Docker) Stop(id string) DockerResult {
    log.Printf("Stopping container %s\n", id)
    ctx := context.Background()

    // Stop the container
    err := d.Client.ContainerStop(
        ctx,
        id,
        container.StopOptions{},
    )
    if err != nil {
        log.Printf("Error while stopping container %s: %v\n", id, err)
        return DockerResult{Error: err,}
    }

    // Remove the container
    err = d.Client.ContainerRemove(
        ctx,
        id,
        types.ContainerRemoveOptions{
            RemoveVolumes: true,
            RemoveLinks: false,
            Force: false,
        },
    )
    if err != nil {
        log.Printf("Error while removing container %s: %v\n", id, err)
        return DockerResult{Error: err,}
    }
    return DockerResult{
        ContainerId: id,
        Action: "stop",
        Result: "success",
        Error: nil,
    }
}


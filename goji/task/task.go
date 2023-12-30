package task

import (
    "github.com/google/uuid"
    "github.com/docker/go-connections/nat"
    "time"
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

package main

import (
    "fmt"
    "os"
    "strings"
    "time"

	"github.com/docker/docker/client"
    "github.com/google/uuid"
    "github.com/golang-collections/collections/queue"

    "goji/manager"
    "goji/node"
    "goji/worker"
    "goji/task"
)


// The code for goji's orchestrator is organized into separate sub-directories inside our project: 
// - manager
// - node
// - scheduler
// - task
// - worker.
//
// A task can be in one of five states: Pending, Scheduled, Running, Completed, or Failed. 
// The worker and manager will use these states to perform actions on tasks, such as stopping and starting them.
// We use interface to allow us to implement multiple schedulers, each with slightly different behavior.

func main() {
    line := strings.Repeat("-", 50) + "\n"
    asterisks := strings.Repeat("*", 25) + "\n"
    fmt.Printf(line)

    // -------------------------------------
    t := task.Task{
        Id: uuid.New(),
        Name: "task-1",
        State: task.Pending,
        Image: "image-1",
        Memory: 1024,
        Disk: 1,
	}

    te := task.TaskEvent{
        Id: uuid.New(),
        State: task.Pending,
        Timestamp: time.Now(),
        Task: t,
	}
    fmt.Printf("task %v\n", t)
    fmt.Printf("task event: %v\n", te)
    // -------------------------------------

    // -------------------------------------
    fmt.Printf(asterisks)
    
    w := worker.Worker{
        Name: "worker-1",
        Queue: *queue.New(),
        Db: make(map[uuid.UUID]*task.Task),
	}
    fmt.Printf("worker %v\n", w)
    w.CollectStats()
    w.RunTask()
    w.StartTask()
    w.StopTask()
    // -------------------------------------

    // -------------------------------------
    fmt.Printf(asterisks)

    m := manager.Manager{
        Pending: *queue.New(),
        TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}
    fmt.Printf("manager: %v\n", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()
    // -------------------------------------

    // -------------------------------------
    fmt.Printf(asterisks)

	n := node.Node{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}
	fmt.Printf("node: %v\n", n)
    // -------------------------------------

    // -------------------------------------
    fmt.Printf(asterisks)

    fmt.Printf("creating a test container\n")
    dockerTask, createResult := createContainer()
    if createResult.Error != nil {
        fmt.Printf("%v\n", createResult.Error)
        os.Exit(1)
    }

    time.Sleep(5 * time.Second)
    fmt.Printf("stopping container %s\n", createResult.ContainerId)
    _ = stopContainer(dockerTask, createResult.ContainerId)

    // -------------------------------------
    fmt.Printf(line)
}

func createContainer() (*task.Docker, *task.DockerResult) {
    config := task.Config{
        Name: "test-container-1",
        Image: "postgres:13",
        Env: []string{
            "POSTGRES_USER=goji", 
            "POSTGRES_PASSWORD=sec",
        },
    }
    docker_client, _ := client.NewClientWithOpts(client.FromEnv)
    d := task.Docker{
        Client: docker_client,
        Config: config,
    }

    result := d.Run()
    if result.Error != nil {
        fmt.Printf("Error: %v\n", result.Error)
        return nil, nil
    }

    fmt.Printf("Container %s is running with config %v\n", result.ContainerId, config)
    return &d, &result
}


func stopContainer(d *task.Docker, id string) *task.DockerResult {
    result := d.Stop(id) 
    if result.Error != nil {
        fmt.Printf("Error: %v\n", result.Error)
        return nil
    }

    fmt.Printf("Container %s has been stopped and removed\n", result.ContainerId)
    return &result
}

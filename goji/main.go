package main

import (
    "fmt"
    "strings"
    "time"
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
    fmt.Printf(line)
}

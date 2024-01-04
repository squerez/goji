package worker

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/golang-collections/collections/queue"

    "cube/task"
)

// The worker requirements are:
// - run tasks as docker containers
// - accept tasks to run from a manager
// - provide relevant statistics to the manager for the purpose of scheduling tasks
// - keep track of its tasks and their state
//
// To accept tasks to run from manager, worker will need a Queue field to ensure tasks are handled in FIFO order
// To keep track of tasks, we will use a field called Db, which will be a map of UUIDs to tasks
// To keep track of all tasks a worker has at any given time, we use TaskCount
type Worker struct {
    Name string
    Queue queue.Queue
    Db map[uuid.UUID]*task.Task
    TaskCount int
}


// To handle the running of a task on the machine where the worker is running, we create a func called RunTask
// It will be responsible for identifying the task's current state, and then either starting or stopping a task based on the state
// Similarly, we will create some func to StartTask and to StopTask
func (w *Worker) RunTask() {
    fmt.Println("I will start or stop a task")
}

func (w *Worker) StartTask() {
    fmt.Println("I will start a task")
}

func (w *Worker) StopTask() {
    fmt.Println("I will stop a task")
}

// To handle the collection of periodic statistics about the weorker, we create a CollectStats func
func (w *Worker) CollectStats() {
    fmt.Println("I will collect stats")
}

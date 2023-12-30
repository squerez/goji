package manager

import (
    "fmt"
    "github.com/google/uuid"
    "github.com/golang-collections/collections/queue"

    "goji/task"
)

// The requirements for the manager are:
// - accept requests from users to start and stop tasks
// - schedule tasks onto worker machines
// - keep track of tasks, states and the machine on which they run
//
// For this, we will create the Manager struct, which will have a queue, represent by the Pending field
// in which tasks will be place upon first being submitted.
// The queue will allow the manager to handle tasks on a FIFO order
//
// The manager will also have two in-mem databases: one to store tasks, one to store task events
// Our manager will also keep track of the workers in the cluster, for this, we will use a field called Workers
// 
// Our manager will also use a field called WorkerTaskMap, which will map strings to task UUIDS 
// and will use a TaskWorkerMap field to find the worker running a task, given tis friendly name, which is a map of UUIDS to strings, where the string is the name of the worker
type Manager struct {
    Pending       queue.Queue
    TaskDb        map[string][]*task.Task
    EventDb       map[string][]*task.TaskEvent
    Workers       []string
    WorkerTaskMap map[string][]uuid.UUID
    TaskWorkerMap map[uuid.UUID]string
}

// The manager needs to schedule tasks onto workers, for this we need a method that selects workers to perform tasks, called SelectWorker
// The method will be responsivle for looking at the requirements specified in Task and evaluating resources available in the pool of workers to see which worker is best suited to run the task.
func (m *Manager) SelectWorker() {
    fmt.Println("I will select an appropriate worker")
}

// To meet the requirement of keeping track of tasks, states and the machined on which they run, we need to 
// create a method called UpdateTasks, which will end up triggering a call to a workers CollectStats method
func (m *Manager) UpdateTasks() {
    fmt.Println("I will update tasks")
}


// To effectively send tasks to workers, we need a method for this
func (m *Manager) SendWork() {
    fmt.Println("I will send work to workers")
}

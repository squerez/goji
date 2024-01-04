package node

// While the Worker is the component that deals with our logical workload, aka tasks, it also has a physical aspect to it.
// The Worker needs to runs on a physical machine itself and causes tasks to be run also on a physical machine.
// For this, it needs more information about the underlying machine, such as stats, to inform the manager for its scheduling decisions.
// This physical aspect of a Worker is called Node.
//
// A Node is a object that represents any machine in our cluster - eg a managers is a type of node, a worker is another type of node.
// To represent its role, we use the Role field.
//
// A Node should have a Name field to be identified in a friendly, an Ip address a manager needs to know to send tasks to it.
//
// Also, it should have physical properties, such as Cores, Memory and Disk space to be used in tasks, represented as maximum ammounts.
// At any point in time, the tasks will be using some ammount of disk and memory, represented by MemoryAllocated and DiskAllocated fields.
// 
// A Node will also have zero or more tasks, tracked via TaskCount field. 
type Node struct {
    Name string
    Ip string
    Cores int
    Memory int
    MemoryAllocated int
    Disk int
    DiskAllocated int
    Role string
    TaskCount int
}

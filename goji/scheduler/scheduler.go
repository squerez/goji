package scheduler 

// The scheduler's requirements are:
// - determine a set of candidate workers on which a task could run
// - score the candidate workers from best to worse
// - pick the worker with the best score
//
// Scheduler will be implemented using interfaces. 
// An interface in go is a contract that specifies a set of behaviours.
// Any type that implements the behaviours can be used anywhere that the interface type is specified
//
// There could be a myriad of ways to select the next worker to fullfill a task
// eg. by spreading the load equally across workers, or filling up on worker completely before moving on to the other
// So using an interface allows the flexibility to change the strategy for choosing the next worker
type Scheduler interface {
    SelectCandidateNodes()
    Score()
    Pick()
} 

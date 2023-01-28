package task

// State represents the state in the run cycle of task
type State string

// State of tasks in run cycle
const (
	// state of newly created task or unfreezed task
	StateIdle State = "idle"
	// task is staged to track
	StateStaged State = "staged"
	// task is queued for checks
	StateQueued State = "queued"
	// task is ready for execution
	StateReady State = "ready"
	// task is waiting for its start condition to satisfy
	StateWaiting State = "waiting"
	// task is running
	StateRunning State = "running"
	// task is successfully executed
	StateSuccess State = "success"
	// task is failed with error
	StateFailure State = "failure"
	// task is aborted forcefully while running
	StateAborted State = "aborted"
	// task is ignored for scheduling
	StateFrozen State = "frozen"
)

func (s State) IsIdle() bool {
	return s == StateIdle
}

func (s State) IsStaged() bool {
	return s == StateStaged
}

func (s State) IsQueued() bool {
	return s == StateQueued
}

func (s State) IsReady() bool {
	return s == StateReady
}

func (s State) IsWaiting() bool {
	return s == StateWaiting
}

func (s State) IsRunning() bool {
	return s == StateRunning
}

func (s State) IsSuccess() bool {
	return s == StateSuccess
}

func (s State) IsFailure() bool {
	return s == StateFailure
}

func (s State) IsAborted() bool {
	return s == StateAborted
}

func (s State) IsFrozen() bool {
	return s == StateFrozen
}

// IsStable checks if state is schedulable
func IsStable(s State) bool {
	return s.IsAborted() || s.IsFailure() || s.IsIdle() || s.IsSuccess()
}

// IsRunnable checks if it can be triggered to run
func IsRunnable(s State) bool {
	return !s.IsQueued() || !s.IsReady() || !s.IsStaged() || !s.IsRunning()
}

// IsTriggered checks if task is staged or set for execution
func IsTriggered(s State) bool {
	return s.IsStaged() || s.IsQueued() || s.IsReady() || s.IsWaiting() || s.IsRunning()
}

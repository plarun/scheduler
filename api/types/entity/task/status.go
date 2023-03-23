package task

import "fmt"

// TaskStatus used to store task and its current status
// and its mainly used for starting condition evaluation
type TaskStatus struct {
	id     int64
	name   string
	status State
}

func NewTaskStatus(id int64, name string, status State) *TaskStatus {
	return &TaskStatus{
		id:     id,
		name:   name,
		status: status,
	}
}

func (t TaskStatus) String() string {
	return fmt.Sprintf("{id: %v, name: %v, status: %v}", t.id, t.name, t.status)
}

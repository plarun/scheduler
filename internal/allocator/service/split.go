package service

import (
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/query"
)

type TaskSplitter struct{}

func NewTaskSplitter() *TaskSplitter {
	return &TaskSplitter{}
}

// Split checks the queued tasks and route them to ready or wait
// based on their start condition status
func (t *TaskSplitter) Split() error {
	// lock the queued tasks for start condition check
	if err := query.LockForConditionCheck(); err != nil {
		return fmt.Errorf("Split: %w", err)
	}

	// pick the locked tasks to in memory
	tasks, err := query.PickQueueLockedTasks()
	if err != nil {
		return fmt.Errorf("Split: %w", err)
	}

	// check start_condition for each task taken into in memory
	// then update the lock flag
	for _, taskid := range tasks {
		ready, err := CheckStartCondition(taskid)
		if err != nil {
			return fmt.Errorf("Split: %w", err)
		}

		if ready {
			if err := query.MoveQueueToReady(taskid); err != nil {
				return fmt.Errorf("Split: %w", err)
			}
		} else {
			if err := query.MoveQueueToWait(taskid); err != nil {
				return fmt.Errorf("Split: %w", err)
			}
		}
	}
	return nil
}

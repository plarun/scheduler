package service

import (
	"fmt"
	"time"

	"github.com/plarun/scheduler/internal/allocator/db/query"
)

type TaskPoller struct {
	cycle time.Duration
}

func NewTaskPoller(cycle time.Duration) *TaskPoller {
	return &TaskPoller{cycle}
}

// Stage performs staging the scheduled tasks
func (t *TaskPoller) Stage() error {
	// lock tasks for staging
	if err := query.LockForStaging(); err != nil {
		return fmt.Errorf("Stage: %w", err)
	}
	// push locked tasks into staging area
	if err := query.StageLockedTasks(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}
	// set task status to "staged"
	if err := query.MarkAsStaged(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}
	// set task flag as staged
	if err := query.SetStagedFlag(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}

	return nil
}

// Poll performs queuing the staged tasks into queue
func (t *TaskPoller) Poll() error {
	// lock staged tasks for queuing
	if err := query.LockForEnqueue(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	// push staged tasks which are locked into queue
	if err := query.EnqueueTasks(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	// mark staged bundle tasks for queuing its tasks (flag=4)
	if err := query.ChangeStagedBundleLock(2, 4); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	// mark tasks of locked bundle for staging
	if err := query.LockBundledTasksForStaging(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	// mark staged bundle tasks after queueing its tasks (flag=5)
	if err := query.ChangeStagedBundleLock(4, 5); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	// set bundle task state to running when flag=5
	if err := query.MarkStagedBundleAsRunning(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	// set status of queued tasks to 'queued'
	if err := query.SetQueueStatus(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	if err := query.SetQueuedFlag(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	return nil
}

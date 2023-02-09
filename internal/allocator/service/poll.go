package service

import (
	"fmt"
	"log"
	"time"

	db "github.com/plarun/scheduler/internal/allocator/db/mysql/query"
)

type TaskPoller struct {
	cycle time.Duration
}

func NewTaskPoller(cycle time.Duration) *TaskPoller {
	return &TaskPoller{cycle}
}

// Stage performs staging the scheduled tasks
func (t *TaskPoller) Stage() error {
	log.Println("Staging...")
	// lock tasks for staging
	if err := db.LockForStaging(); err != nil {
		return fmt.Errorf("Stage: %w", err)
	}
	// push locked tasks into staging area
	if err := db.StageLockedTasks(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}
	// set task status to "staged"
	if err := db.MarkAsStaged(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}
	// set task flag as staged
	if err := db.SetStagedFlag(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}

	return nil
}

// Poll performs queuing the staged tasks into queue
func (t *TaskPoller) Poll() error {
	log.Println("Polling...")

	log.Println("LockForEnqueue")
	if err := db.LockForEnqueue(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	log.Println("EnqueueTasks")
	if err := db.EnqueueTasks(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	log.Println("LockStagedBundles")
	if err := db.LockStagedBundles(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	log.Println("StageBundledTasks")
	if err := db.StageBundledTasks(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	log.Println("SetQueueStatus")
	if err := db.SetQueueStatus(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	log.Println("SetQueuedFlag")
	if err := db.SetQueuedFlag(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	return nil
}

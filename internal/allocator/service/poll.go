package service

import (
	"fmt"
	"time"

	db "github.com/plarun/scheduler/internal/allocator/db/mysql/query"
)

type TaskPoller struct {
	cycle time.Duration
}

func NewTaskPoller(cycle time.Duration) *TaskPoller {
	return &TaskPoller{cycle}
}

func (t *TaskPoller) Stage() error {

	if err := db.StageTasks(int(t.cycle)); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}

	if err := db.SetStagedStatus(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}

	if err := db.MigrateStage(); err != nil {
		return fmt.Errorf("Stage: %v", err)
	}

	return nil
}

func (t *TaskPoller) Poll() error {

	if err := db.SetFlagPreQueue(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	if err := db.QueueTasks(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	if err := db.SetStatusQueue(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	if err := db.SetFlagPostQueue(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	return nil
}

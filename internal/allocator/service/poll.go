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

func (t *TaskPoller) Poll() error {

	if err := db.StageTasks(int(t.cycle)); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	if err := db.SetStagedStatus(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	if err := db.MigrateStage(); err != nil {
		return fmt.Errorf("Poll: %v", err)
	}

	return nil
}

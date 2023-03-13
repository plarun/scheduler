package service

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/eventserver/db/query"
)

func ChangeTaskState(ctx context.Context, id int64, state task.State) error {
	// change the task state
	if err := query.SetTaskStatus(id, state); err != nil {
		return fmt.Errorf("ChangeTaskState: %w", err)
	}
	return nil
}

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

	// invoke the dependent tasks if any of them are waiting
	if state.IsSuccess() || state.IsFailure() {
		// get all the dependent task ids of this task
		// pass them to check on sched_wait
		// then awake them if any

		// rpc to allocator to awake dep tasks from waiting
	}
	return nil
}

package service

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/api/types/condition"
	"github.com/plarun/scheduler/internal/eventserver/db/query"
)

func AwakeWaitingDependentTasks(ctx context.Context, id int64) error {
	var dep []int64

	// check if start condition satisfied
	condStr, err := query.GetStartCondition(id)
	if err != nil {
		return fmt.Errorf("AwakeWaitingDependentTasks: %w", err)
	}

	// build task's start condition string into expression to evaluate
	expr, err := condition.Build(condStr)
	if err != nil {
		return fmt.Errorf("AwakeWaitingDependentTasks: %w", err)
	}

	// get status of dependent tasks'status
	depStatus, err := query.GetDependentTasksStatus(id)
	if err != nil {
		return fmt.Errorf("AwakeWaitingDependentTasks: %w", err)
	}

	if len(depStatus) == 0 {
		return nil
	}

	if ok, err := expr.Check(depStatus); err != nil {
		return fmt.Errorf("AwakeWaitingDependentTasks: %w", err)
	}

	depTasks := make([]string, 0)

	// move the dependent tasks into ready queue from waiting queue
	if err := query.MoveWaitToReady(depStatus); err != nil {
		return fmt.Errorf("AwakeWaitingDependentTasks: %w", err)
	}

	return nil
}

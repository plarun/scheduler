package service

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	db "github.com/plarun/scheduler/internal/allocator/db/mysql/query"
)

// Init loads the ConditionChecker which are required for evaluation
func CheckStartCondition(id int64) (bool, error) {
	// get start condition of task
	cond, err := db.GetStartCondition(id)
	if err != nil {
		return false, fmt.Errorf("ConditionChecker.Init: %w", err)
	}

	// get current status of distinct tasks in start condition
	stat, err := db.GetDependentTasksStatus(id)
	if err != nil {
		return false, fmt.Errorf("ConditionChecker.Init: %w", err)
	}

	cc := task.NewConditionChecker(id, cond, stat)
	if ok, err := cc.Check(); err != nil {
		return false, fmt.Errorf("ConditionChecker.Init: %w", err)
	} else {
		return ok, nil
	}
}

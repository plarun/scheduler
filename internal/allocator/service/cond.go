package service

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/condition"
	"github.com/plarun/scheduler/api/types/entity/task"
	db "github.com/plarun/scheduler/internal/allocator/db/mysql/query"
)

type ConditionChecker struct {
	condTaskStatus map[string]task.State
	task           string
	cond           string
	expr           condition.Expression
	initiated      bool
}

func NewConditionChecker(name string) *ConditionChecker {
	return &ConditionChecker{
		task:      name,
		initiated: false,
	}
}

func (c *ConditionChecker) Init() error {
	// get start condition of task
	if cond, err := db.GetStartCondition(c.task); err != nil {
		return fmt.Errorf("ConditionChecker.Init: %w", err)
	} else {
		c.cond = cond
	}

	// get current status of distinct tasks in start condition
	if stat, err := db.GetPrerequisitesTaskStatus(c.task); err != nil {
		return fmt.Errorf("ConditionChecker.Init: %w", err)
	} else {
		for tsk, status := range stat {
			c.condTaskStatus[tsk] = task.State(status)
		}
	}

	// build start condition expression for evaluation
	if expr, err := condition.Build(c.cond); err != nil {
		return fmt.Errorf("ConditionChecker.Init: %w", err)
	} else {
		c.expr = expr
	}

	c.initiated = true
	return nil
}

func (c *ConditionChecker) Check() (bool, error) {
	if !c.initiated {
		return false, fmt.Errorf("ConditionChecker.Check: Initiation required")
	}
	if res, err := c.eval(); err != nil {
		return false, fmt.Errorf("ConditionChecker.Check: %w", err)
	} else {
		return res, nil
	}
}

func (c *ConditionChecker) eval() (bool, error) {
	return false, nil
}

package service

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/condition"
	"github.com/plarun/scheduler/api/types/entity/task"
	db "github.com/plarun/scheduler/internal/allocator/db/mysql/query"
)

// ConditionChecker represents a type to check whether start condition
// of task is satisfied or not
type ConditionChecker struct {
	condTaskStatus map[string]task.State
	taskid         int
	cond           string
	expr           condition.Expression
	initiated      bool
}

func NewConditionChecker(id int) *ConditionChecker {
	return &ConditionChecker{
		taskid:    id,
		initiated: false,
	}
}

// Init loads the ConditionChecker which are required for evaluation
func (c *ConditionChecker) Init() error {
	// get start condition of task
	if cond, err := db.GetStartCondition(c.taskid); err != nil {
		return fmt.Errorf("ConditionChecker.Init: %w", err)
	} else {
		c.cond = cond
	}

	// get current status of distinct tasks in start condition
	if stat, err := db.GetPrerequisitesTaskStatus(c.taskid); err != nil {
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

// Check checks whether a start condition is satisfied or not
func (c *ConditionChecker) Check() (bool, error) {
	if !c.initiated {
		if err := c.Init(); err != nil {
			return false, err
		}
	}
	if res, err := c.eval(c.expr); err != nil {
		return false, fmt.Errorf("ConditionChecker.Check: %w", err)
	} else {
		return res, nil
	}
}

// eval evaluates the start condition as boolean expression
func (c *ConditionChecker) eval(cur condition.Expression) (bool, error) {
	if cur.IsWrapper() {
		wrapper := cur.(*condition.Wrapper)
		// evaluate each sub expressions
		for i := 0; i < len(wrapper.Conditions); i++ {
			cond := wrapper.Conditions[i]
			res, err := c.eval(cond)
			if err != nil {
				return false, err
			}
			cond.SetResult(res)
		}
		// evaluate the result of wrapper
		// here cond can be either clause or wrapper
		op := condition.OperatorEmpty
		for i := 0; i < len(wrapper.Conditions); i++ {
			cond := wrapper.Conditions[i]
			if i == 0 {
				wrapper.SetResult(cond.GetResult())
				op = cond.GetOperator()
			} else {
				res := wrapper.GetResult()
				switch op {
				case condition.OperatorAnd:
					res = res && cond.GetResult()
				case condition.OperatorOr:
					res = res || cond.GetResult()
				case condition.OperatorEmpty:
					if i+1 != len(wrapper.Conditions) {
						return false, fmt.Errorf("ConditionChecker.eval: empty operator")
					}
				}
				wrapper.SetResult(res)
				op = cond.GetOperator()
			}
		}
		return wrapper.GetResult(), nil
	} else {
		clause := cur.(*condition.Clause)
		state := c.condTaskStatus[clause.TaskName]
		clause.SetResult(evalClause(clause.Status, state))
		return clause.GetResult(), nil
	}
}

// evalClause is helper function to convert state type in start condition
// into corresponding state type in task
func evalClause(st string, state task.State) bool {
	// ignoring the frozen task
	if state.IsFrozen() {
		return true
	}
	switch st {
	case "su":
		if state.IsSuccess() {
			return true
		}
	case "fa":
		if state.IsFailure() {
			return true
		}
	case "nr":
		if task.IsStable(state) {
			return true
		}
	}
	return false
}

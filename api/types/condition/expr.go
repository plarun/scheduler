package condition

import (
	"fmt"
	"regexp"
	"strings"
)

type Operator string

const (
	OperatorEmpty Operator = ""
	OperatorAnd   Operator = "&"
	OperatorOr    Operator = "|"
)

func stringToCondOperator(str string) Operator {
	switch str {
	case "&":
		return OperatorAnd
	case "|":
		return OperatorOr
	default:
		return OperatorEmpty
	}
}

func condOperatorToString(op Operator) string {
	switch op {
	case OperatorAnd:
		return "&"
	case OperatorOr:
		return "|"
	default:
		return ""
	}
}

type Expression interface {
	build() string
	IsWrapper() bool
	addChild(Expression)
	setParent(*Wrapper)
	getParent() *Wrapper
	String() string
	SetResult(bool)
	GetResult() bool
	GetOperator() Operator
}

// addWrapper creates and adds a new wrapper to current condition
// parent is wrapper and child is also wrapper
func addWrapper(parent Expression) (*Wrapper, error) {
	if !parent.IsWrapper() {
		return nil, fmt.Errorf("addWrapper: parent is not wrapper")
	}

	child := newWrapper()
	parent.addChild(child)
	return child, nil
}

// addCondition creates and adds a task condition to current condition
// parent is wrapper and child is condclause
func addCond(parent Expression, status, tsk, op string) (*Clause, error) {
	if !parent.IsWrapper() {
		return nil, fmt.Errorf("addCond: parent is not wrapper")
	}

	child := newClause(status, tsk, op)
	parent.addChild(child)
	return child, nil
}

// getDistinctJobs returns list of distinct tasks in task's start condition
func GetDistinctTasks(condStr string) []string {
	expr, _ := Build(condStr)
	set := make(map[string]bool)

	var que []Expression
	que = append(que, expr)

	for len(que) != 0 {
		curr := que[0]
		que = que[1:]
		if curr.IsWrapper() {
			w := curr.(*Wrapper)
			que = append(que, w.Conditions...)
		} else {
			c := curr.(*Clause)
			set[c.TaskName] = true
		}
	}

	var jobs []string
	for job := range set {
		jobs = append(jobs, job)
	}

	return jobs
}

// buildCondition parses the given string representation of task
// condition and builds corresponding tree format (ConditionClause)
func Build(condition string) (Expression, error) {
	condition = strings.ReplaceAll(condition, " ", "")
	size := len(condition)
	isNewCondition, isTaskName, closeWrap, mayBeOperator := true, false, false, false

	var status, taskName string
	var root, curr *Wrapper
	root = newWrapper()
	curr = root

	jobRegex, _ := regexp.Compile("^[0-9a-zA-Z_]$")

	for i := 0; i < size; i++ {

		// wrap
		if condition[i:i+1] == "(" {
			w, err := addWrapper(curr)
			if err != nil {
				return nil, err
			}
			curr = w
			continue
		}

		// task condition
		if isNewCondition {
			// status
			if i+2 < size {
				tag := strings.ToLower(condition[i : i+2])
				if tag != "su" && tag != "fa" && tag != "nr" {
					return nil, fmt.Errorf("BuildCondition: invalid condition tag")
				}
				status = tag
				i += 2
			}

			// syntax check
			if i >= size || condition[i:i+1] != "(" {
				return nil, fmt.Errorf("BuildCondition: condition syntax expecting (")
			}
			i += 1

			// syntax check
			if i >= size || condition[i:i+1] == ")" {
				return nil, fmt.Errorf("BuildCondition: condition has empty task")
			}
			isTaskName = true
		}

		// task name in start condition
		if isTaskName {
			var j int
			for j = i; condition[i:i+1] != ")"; i++ {
				if ok := jobRegex.MatchString(condition[i : i+1]); !ok {
					return nil, fmt.Errorf("BuildCondition: invalid expression")
				}
				if i+1 >= size {
					return nil, fmt.Errorf("BuildCondition: condition syntax expecting )")
				}
			}
			taskName = condition[j:i]
			i++
			isTaskName, mayBeOperator = false, true
		}

		// condition operator
		if mayBeOperator {
			if i == size {
				addCond(curr, status, taskName, "")
			} else {
				join := condition[i : i+1]
				if join == "&" || join == "|" {
					// condition clause ending with operator
					if i+1 == size {
						return nil, fmt.Errorf("BuildCondition: incomplete condition clause")
					}
					addCond(curr, status, taskName, join)
				} else {
					addCond(curr, status, taskName, "")
					closeWrap, isNewCondition = true, false
				}
				// reset
				status = ""
				taskName = ""
			}
			isTaskName, mayBeOperator = false, false
		}

		// closeWrap
		if closeWrap {
			if root == curr || i >= size || condition[i:i+1] != ")" {
				return nil, fmt.Errorf("BuildCondition: condition syntax expecting )")
			}

			// join of group
			if i+1 < size && (condition[i+1:i+2] == "&" || condition[i+1:i+2] == "|") {
				i++
				curr.Operator = stringToCondOperator(condition[i : i+1])
				isNewCondition = true
				closeWrap = false
			}

			curr = curr.getParent()
		}
	}

	if curr != root {
		return nil, fmt.Errorf("BuildCondition: incomplete condition clause")
	}
	return root, nil
}

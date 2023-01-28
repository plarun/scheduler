package service

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/condition"
	"github.com/plarun/scheduler/api/types/entity/task"
	db "github.com/plarun/scheduler/internal/validator/db/mysql/query"
	"github.com/plarun/scheduler/internal/validator/errors"
	"github.com/plarun/scheduler/proto"
)

type Checker struct {
	insertTasks map[string]*localTask
	updateTasks map[string]*localTask
	deleteTasks map[string]*localTask
}

type localTask struct {
	name     string
	tasktype task.Type
	pos      int
}

func newLocalTask(name string, tp task.Type, pos int) *localTask {
	return &localTask{
		name:     name,
		tasktype: tp,
		pos:      pos,
	}
}

func (lt *localTask) String() string {
	return fmt.Sprintf("{task_name: %s, task_type: %s, task_pos: %d}", lt.name, lt.tasktype, lt.pos)
}

func NewChecker() *Checker {
	return &Checker{
		insertTasks: make(map[string]*localTask),
		updateTasks: make(map[string]*localTask),
		deleteTasks: make(map[string]*localTask),
	}
}

func (chk *Checker) String() string {
	var insTasks, updTasks, delTasks []string

	for _, t := range chk.insertTasks {
		insTasks = append(insTasks, t.String())
	}
	for _, t := range chk.updateTasks {
		updTasks = append(updTasks, t.String())
	}
	for _, t := range chk.deleteTasks {
		delTasks = append(delTasks, t.String())
	}

	return fmt.Sprintf(
		"{insert_tasks: %v, update_tasks: %v, delete_tasks: %v}",
		insTasks, updTasks, delTasks,
	)
}

func (chk *Checker) CheckExistance(tsks []*proto.ValidatedTaskEntity) error {
	for i, tsk := range tsks {
		if tsk == nil {
			log.Fatal("nil task")
		}
		// check task existance based on action
		if err := chk.checkTask(tsk, i); err != nil {
			return fmt.Errorf("CheckExistance: %w", err)
		}

		// parent task should exists in db or should be inserted before this
		if err := chk.checkParent(tsk, i); err != nil {
			return fmt.Errorf("CheckExistance: %w", err)
		}

		// condition tasks should exists in db or should be available
		// before this task in the JIL
		if err := chk.checkPredecessors(tsk); err != nil {
			return fmt.Errorf("CheckExistance: %w", err)
		}
	}

	return nil
}

func (chk *Checker) checkTask(tsk *proto.ValidatedTaskEntity, pos int) error {
	switch tsk.Action {
	case string(task.ActionInsert):
		return chk.checkInsertTask(tsk, pos)
	case string(task.ActionUpdate):
		return chk.checkUpdateTask(tsk, pos)
	case string(task.ActionDelete):
		return chk.checkDeleteTask(tsk, pos)
	}

	return nil
}

func (chk *Checker) checkInsertTask(tsk *proto.ValidatedTaskEntity, pos int) error {
	if exists, err := db.TaskExists(tsk.Name); err != nil {
		return fmt.Errorf("checkInsertTask: %w", err)
	} else if exists {
		return errors.ErrTaskAlreadyExist
	}

	chk.insertTasks[tsk.Name] = newLocalTask(tsk.Name, task.Type(tsk.Type.Value), pos+1)
	return nil
}

func (chk *Checker) checkUpdateTask(tsk *proto.ValidatedTaskEntity, pos int) error {
	if exists, err := db.TaskExists(tsk.Name); err != nil {
		return fmt.Errorf("CheckExistance: %w", err)
	} else if !exists {
		return errors.ErrTaskNotExist
	}

	chk.updateTasks[tsk.Name] = newLocalTask(tsk.Name, task.Type(tsk.Type.Value), pos+1)
	return nil
}

func (chk *Checker) checkDeleteTask(tsk *proto.ValidatedTaskEntity, pos int) error {
	if exists, err := db.TaskExists(tsk.Name); err != nil {
		return fmt.Errorf("CheckExistance: %w", err)
	} else if !exists {
		return errors.ErrTaskNotExist
	}

	chk.deleteTasks[tsk.Name] = newLocalTask(tsk.Name, task.Type(tsk.Type.Value), pos+1)
	return nil
}

func (chk *Checker) checkParent(tsk *proto.ValidatedTaskEntity, pos int) error {
	if tsk.Parent.Flag == proto.NullableFlag_Available {
		// check bundle task in local
		if lTsk, exists := chk.insertTasks[tsk.Parent.Value]; exists && lTsk.pos < pos+1 {
			if !lTsk.tasktype.IsBundle() {
				return fmt.Errorf("checkParent: task %s is not bundle type, it cannot be parent", lTsk.name)
			}
		} else {
			// check bundle task in db
			tasktype, err := db.GetTaskType(tsk.Parent.Value)
			if err != nil {
				return fmt.Errorf("checkParent: %v", err)
			} else if task.Type(tasktype).IsBundle() {
				return fmt.Errorf("checkParent: parent task should be bundle type")
			}
		}
	}
	return nil
}

func (chk *Checker) checkPredecessors(tsk *proto.ValidatedTaskEntity) error {
	var pred []string

	if tsk.Condition.Flag == proto.NullableFlag_Available {
		pred = condition.GetDistinctTasks(tsk.Condition.Value)
	}

	for _, condTask := range pred {
		lTsk := chk.insertTasks[condTask]

		if lTsk == nil { // task doesn't exist in current JIL
			if exists, err := db.TaskExists(condTask); err != nil {
				return fmt.Errorf("CheckExistance: %v", err)
			} else if !exists {
				return errors.ErrTaskNotExist
			}
		} else if tsk.Name == condTask { // self reference
			return fmt.Errorf("self reference is not allowed")
		}
	}
	return nil
}

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
	insertJobs map[string]*localTask
	updateJobs map[string]*localTask
	deleteJobs map[string]*localTask
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

func (lJob *localTask) String() string {
	return fmt.Sprintf("{job_name: %s, job_type: %s, job_pos: %d}", lJob.name, lJob.tasktype, lJob.pos)
}

func NewChecker() *Checker {
	return &Checker{
		insertJobs: make(map[string]*localTask),
		updateJobs: make(map[string]*localTask),
		deleteJobs: make(map[string]*localTask),
	}
}

func (chk *Checker) String() string {
	var insJobs, updJobs, delJobs, condJobs []string

	for _, insJob := range chk.insertJobs {
		insJobs = append(insJobs, insJob.String())
	}
	for _, updJob := range chk.updateJobs {
		updJobs = append(updJobs, updJob.String())
	}
	for _, delJob := range chk.deleteJobs {
		delJobs = append(delJobs, delJob.String())
	}

	return fmt.Sprintf(
		"{insert_jobs: %v, update_jobs: %v, delete_jobs: %v, condition_jobs: %v}",
		insJobs, updJobs, delJobs, condJobs,
	)
}

func (chk *Checker) CheckExistance(validatedJobs []*proto.ValidatedTaskEntity) error {
	for i, tsk := range validatedJobs {
		if tsk == nil {
			log.Fatal("nil task")
		}
		// check task existance based on action
		if err := chk.checkJob(tsk, i); err != nil {
			return fmt.Errorf("CheckExistance: %w", err)
		}

		// parent task should exists in db or should be inserted before this
		if err := chk.checkParent(tsk, i); err != nil {
			return fmt.Errorf("CheckExistance: %w", err)
		}

		// condition jobs should exists in db or should be available
		// before this task in the JIL
		if err := chk.checkPredecessors(tsk); err != nil {
			return fmt.Errorf("CheckExistance: %w", err)
		}
	}

	return nil
}

func (chk *Checker) checkJob(tsk *proto.ValidatedTaskEntity, pos int) error {
	switch tsk.Action {
	case string(task.ActionInsert):
		return chk.checkInsertJob(tsk, pos)
	case string(task.ActionUpdate):
		return chk.checkUpdateJob(tsk, pos)
	case string(task.ActionDelete):
		return chk.checkDeleteJob(tsk, pos)
	}

	return nil
}

func (chk *Checker) checkInsertJob(tsk *proto.ValidatedTaskEntity, pos int) error {
	if exists, err := db.JobExists(tsk.Name); err != nil {
		return fmt.Errorf("checkInsertJob: %w", err)
	} else if exists {
		return errors.ErrJobAlreadyExist
	}

	chk.insertJobs[tsk.Name] = newLocalTask(tsk.Name, task.Type(tsk.Type.Value), pos+1)
	return nil
}

func (chk *Checker) checkUpdateJob(tsk *proto.ValidatedTaskEntity, pos int) error {
	if exists, err := db.JobExists(tsk.Name); err != nil {
		return fmt.Errorf("CheckExistance: %w", err)
	} else if !exists {
		return errors.ErrJobNotExist
	}

	chk.updateJobs[tsk.Name] = newLocalTask(tsk.Name, task.Type(tsk.Type.Value), pos+1)
	return nil
}

func (chk *Checker) checkDeleteJob(tsk *proto.ValidatedTaskEntity, pos int) error {
	if exists, err := db.JobExists(tsk.Name); err != nil {
		return fmt.Errorf("CheckExistance: %w", err)
	} else if !exists {
		return errors.ErrJobNotExist
	}

	chk.deleteJobs[tsk.Name] = newLocalTask(tsk.Name, task.Type(tsk.Type.Value), pos+1)
	return nil
}

func (chk *Checker) checkParent(tsk *proto.ValidatedTaskEntity, pos int) error {
	if tsk.Parent.Flag == proto.NullableFlag_Available {
		// check bundle task in local
		if lTsk, exists := chk.insertJobs[tsk.Parent.Value]; exists && lTsk.pos < pos+1 {
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

	for _, condJob := range pred {
		lTsk := chk.insertJobs[condJob]

		if lTsk == nil { // task doesn't exist in current JIL
			if exists, err := db.JobExists(condJob); err != nil {
				return fmt.Errorf("CheckExistance: %v", err)
			} else if !exists {
				return errors.ErrJobNotExist
			}
		} else if tsk.Name == condJob { // self reference
			return fmt.Errorf("self reference is not allowed")
		}
	}
	return nil
}

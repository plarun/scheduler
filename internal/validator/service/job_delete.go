package service

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	db "github.com/plarun/scheduler/internal/validator/db/mysql/query"
	er "github.com/plarun/scheduler/internal/validator/errors"
	"github.com/plarun/scheduler/proto"
)

type deleteValidation struct {
	input  *proto.ParsedTaskEntity
	output *proto.ValidatedTaskEntity
	task   *task.TaskEntity
}

func newDeleteValidation(tsk *proto.ParsedTaskEntity) *deleteValidation {
	return &deleteValidation{
		input:  tsk,
		output: newValidatedTask(),
	}
}

func (del *deleteValidation) Get() *proto.ValidatedTaskEntity {
	return del.output
}

func (del *deleteValidation) Validate() error {
	badErr := func(err error) error {
		return &er.BadTaskDefError{
			Action: string(task.ActionUpdate),
			Target: del.input.Target,
			Err:    err,
		}
	}

	tsk, otask := del.task, del.output

	// set task action type
	otask.Action = string(task.ActionDelete)

	// validate task name
	if err := checkFieldTaskName(tsk.Name()); err != nil {
		return fmt.Errorf("insertValidation.validate: %w", err)
	}
	otask.Name = tsk.Name()

	// task should already exist
	if exist, err := db.TaskExists(otask.Name); err != nil {
		return err
	} else if !exist {
		return badErr(er.ErrTaskNotExist)
	}

	// other fields should not be available for delete action

	if tsk.HasField(task.FIELD_COMMAND) ||
		tsk.HasField(task.FIELD_CONDITION) ||
		tsk.HasField(task.FIELD_LABEL) ||
		tsk.HasField(task.FIELD_MACHINE) ||
		tsk.HasField(task.FIELD_PARENT) ||
		tsk.HasField(task.FIELD_PRIORITY) ||
		tsk.HasField(task.FIELD_RUN_DAYS) ||
		tsk.HasField(task.FIELD_RUN_WINDOW) ||
		tsk.HasField(task.FIELD_START_MINS) ||
		tsk.HasField(task.FIELD_START_TIMES) ||
		tsk.HasField(task.FIELD_ERR_LOG_FILE) ||
		tsk.HasField(task.FIELD_OUT_LOG_FILE) ||
		tsk.HasField(task.FIELD_TYPE) {
		return badErr(er.ErrBadDeleteTask)
	}

	return nil
}

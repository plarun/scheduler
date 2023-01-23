package service

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	db "github.com/plarun/scheduler/internal/validator/db/mysql/query"
	er "github.com/plarun/scheduler/internal/validator/errors"
	"github.com/plarun/scheduler/proto"
)

type insertValidation struct {
	input  *proto.ParsedTaskEntity
	output *proto.ValidatedTaskEntity
	task   *task.TaskEntity
}

func newInsertValidation(tsk *proto.ParsedTaskEntity) *insertValidation {
	return &insertValidation{
		input:  tsk,
		output: newValidatedTask(),
	}
}

func (ins *insertValidation) Get() *proto.ValidatedTaskEntity {
	return ins.output
}

func (ins *insertValidation) Validate() error {
	// error to point out the error field with its input and job name
	invalidErr := func(err error, field, value string) error {
		return &er.InvalidJobFieldError{
			Action: string(task.ActionInsert),
			Target: ins.input.Target,
			Err:    err,
			Field:  field,
			Value:  value,
		}
	}

	if t, err := prepareTask(ins.input); err != nil {
		return err
	} else {
		ins.task = t
	}

	tsk, otask := ins.task, ins.output

	// set task action type
	otask.Action = string(task.ActionInsert)

	// validate job name
	if err := checkFieldJobName(tsk.Name()); err != nil {
		return fmt.Errorf("insertValidation.validate: %w", err)
	}
	otask.Name = tsk.Name()

	// task should not already exist
	if exist, err := db.JobExists(otask.Name); err != nil {
		return err
	} else if exist {
		return invalidErr(er.ErrJobAlreadyExist, "", otask.Name)
	}

	// validate field 'type'
	// madatory field
	if f, ok := tsk.GetFieldType(); ok {
		otask.Type.Flag = proto.NullableFlag_Available
		otask.Type.Value = string(f.Value())
	} else {
		return invalidErr(er.ErrFieldRequired, f.Name(), "")
	}

	// validate field 'command'
	if f, ok := tsk.GetFieldCommand(); ok {
		// bundle task should not have command
		if typ, _ := tsk.GetType(); typ.IsBundle() {
			return invalidErr(er.ErrFieldNotRequired, f.Name(), f.Value())
		}
		// callable should not have empty command
		if f.Empty() {
			return invalidErr(er.ErrNonEmptyValueRequired, f.Name(), "")
		}
		otask.Command.Flag = proto.NullableFlag_Available
		otask.Command.Value = f.Value()
	} else {
		otask.Command.Flag = proto.NullableFlag_NotAvailable
		// callable task should have command
		if typ, _ := tsk.GetType(); typ.IsCallable() {
			return invalidErr(er.ErrFieldRequired, f.Name(), "")
		}
	}

	// validate field 'condition'
	if f, ok := tsk.GetFieldCondition(); !ok || f.Empty() {
		otask.Condition.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.Condition.Flag = proto.NullableFlag_Available
		otask.Condition.Value = f.Value()
	}

	// validate field 'parent'
	if f, ok := tsk.GetFieldParent(); !ok || f.Empty() {
		otask.Parent.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.Parent.Flag = proto.NullableFlag_Available
		otask.Parent.Value = f.Value()
	}

	// validate field 'machine'
	if f, ok := tsk.GetFieldMachine(); !ok || f.Empty() {
		otask.Machine.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.Machine.Flag = proto.NullableFlag_Available
		otask.Machine.Value = f.Value()
	}

	// validate field 'std_out_log'
	if f, ok := tsk.GetFieldOutLogFile(); !ok || f.Empty() {
		otask.OutLogFile.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.OutLogFile.Flag = proto.NullableFlag_Available
		otask.OutLogFile.Value = f.Value()
	}

	// validate field 'std_err_log'
	if f, ok := tsk.GetFieldErrLogFile(); !ok || f.Empty() {
		otask.ErrLogFile.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.ErrLogFile.Flag = proto.NullableFlag_Available
		otask.ErrLogFile.Value = f.Value()
	}

	// validate field 'label'
	if f, ok := tsk.GetFieldLabel(); !ok || f.Empty() {
		otask.Label.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.Label.Flag = proto.NullableFlag_Available
		otask.Label.Value = f.Value()
	}

	// validate field 'profile'
	if f, ok := tsk.GetFieldProfile(); !ok || f.Empty() {
		otask.Profile.Flag = proto.NullableFlag_NotAvailable
	} else {
		if typ, _ := tsk.GetType(); typ.IsBundle() {
			return invalidErr(er.ErrFieldNotRequired, f.Name(), f.Value())
		}
		otask.Profile.Flag = proto.NullableFlag_Available
		otask.Profile.Value = f.Value()
	}

	// validate field 'priority'
	if f, ok := tsk.GetFieldPriority(); !ok {
		otask.Priority.Flag = proto.NullableFlag_Available
		otask.Priority.Value = 0
	} else {
		otask.Priority.Flag = proto.NullableFlag_Available
		otask.Priority.Value = f.Value()
	}

	// validate field 'run_days'
	if f, ok := tsk.GetFieldRunDays(); !ok {
		otask.RunDays.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.RunDays.Flag = proto.NullableFlag_Available
		otask.RunDays.Value = int32(f.Value())
	}

	// validate field 'start_times'
	if f, ok := tsk.GetFieldStartTimes(); !ok {
		otask.StartTimes.Flag = proto.NullableFlag_NotAvailable
	} else {
		otask.StartTimes.Flag = proto.NullableFlag_Available
		otask.StartTimes.Value = f.Value()
	}

	// validate field 'start_mins'
	if f, ok := tsk.GetFieldStartMins(); !ok {
		otask.StartMins.Flag = proto.NullableFlag_NotAvailable
	} else {
		val := make([]int32, 0)
		for _, m := range f.Value() {
			val = append(val, int32(m))
		}

		otask.StartMins.Flag = proto.NullableFlag_Available
		otask.StartMins.Value = val
	}

	// validate field 'run_window'
	if f, ok := tsk.GetFieldRunWindow(); !ok {
		otask.RunWindow.Flag = proto.NullableFlag_NotAvailable
	} else {
		val := f.Value()
		otask.RunWindow.Flag = proto.NullableFlag_Available
		otask.RunWindow.Value.Start = val[0]
		otask.RunWindow.Value.End = val[1]
	}

	var hasST, hasSM, hasRW bool
	if otask.StartTimes.Flag == proto.NullableFlag_Available {
		hasST = true
	}
	if otask.StartMins.Flag == proto.NullableFlag_Available {
		hasSM = true
	}
	if otask.RunWindow.Flag == proto.NullableFlag_Available {
		hasRW = true
	}

	// batch run should not have window run fields 'start_mins' or 'run_window'
	if hasST && (hasSM || hasRW) {
		return er.ErrBatchWindowRun
	}

	// window run should have both 'start_mins' and 'run_window'
	if hasSM != hasRW {
		return er.ErrIncompleteWindowRun
	}

	return nil
}

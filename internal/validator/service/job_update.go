package service

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	db "github.com/plarun/scheduler/internal/validator/db/mysql/query"
	er "github.com/plarun/scheduler/internal/validator/errors"
	"github.com/plarun/scheduler/proto"
)

type updateValidation struct {
	input  *proto.ParsedTaskEntity
	output *proto.ValidatedTaskEntity
	task   *task.TaskEntity
}

func newUpdateValidation(tsk *proto.ParsedTaskEntity) *updateValidation {
	return &updateValidation{
		input:  tsk,
		output: newValidatedTask(),
	}
}

func (upd *updateValidation) Get() *proto.ValidatedTaskEntity {
	return upd.output
}

func (upd *updateValidation) Validate() error {
	// error to point out the error field with its input and job name
	invalidErr := func(err error, field, value string) error {
		return &er.InvalidJobFieldError{
			Action: string(task.ActionUpdate),
			Target: upd.input.Target,
			Err:    err,
			Field:  field,
			Value:  value,
		}
	}

	badErr := func(err error) error {
		return &er.BadJobDefError{
			Action: string(task.ActionUpdate),
			Target: upd.input.Target,
			Err:    err,
		}
	}

	if t, err := prepareTask(upd.input); err != nil {
		return err
	} else {
		upd.task = t
	}

	tsk, otask := upd.task, upd.output

	// set task action type
	otask.Action = string(task.ActionUpdate)

	// validate task name
	if err := checkFieldJobName(tsk.Name()); err != nil {
		return fmt.Errorf("insertValidation.validate: %w", err)
	}
	otask.Name = tsk.Name()

	// task should already exist
	if exist, err := db.JobExists(otask.Name); err != nil {
		return err
	} else if !exist {
		return badErr(er.ErrJobNotExist)
	}

	// task should already exist for an update and get the run flag of job
	var runFlag task.RunType
	if rf, err := db.GetRunDetails(otask.Name); err != nil {
		return fmt.Errorf("updateValidation.validate: %w", err)
	} else {
		runFlag = convertRunFlag(rf)
	}

	// type of the task
	typ, err := db.GetTaskType(otask.Name)
	if err != nil {
		return fmt.Errorf("updateValidation.validate: %w", err)
	}
	taskType := task.Type(typ)

	// validate field 'type'
	if f, ok := tsk.GetFieldType(); ok {
		return invalidErr(er.ErrNonEditableTypeAttr, f.Name(), string(f.Value()))
	}

	// validate field 'parent'
	if f, ok := tsk.GetFieldParent(); !ok {
		otask.Parent.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.Parent.Flag = proto.NullableFlag_Empty
	} else {
		if err := checkFieldJobName(f.Value()); err != nil {
			return invalidErr(err, f.Name(), f.Value())
		}
		otask.Parent.Flag = proto.NullableFlag_Available
		otask.Parent.Value = f.Value()
	}

	// validate field 'machine'
	if f, ok := tsk.GetFieldMachine(); !ok {
		otask.Machine.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.Machine.Flag = proto.NullableFlag_Empty
	} else {
		otask.Machine.Flag = proto.NullableFlag_Available
		otask.Machine.Value = f.Value()
	}

	// validate field 'command'
	if f, ok := tsk.GetFieldCommand(); !ok {
		otask.Command.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		if taskType.IsCallable() {
			return invalidErr(er.ErrNonEmptyValueRequired, f.Name(), f.Value())
		}
		// ignore empty command for bundle
		otask.Command.Flag = proto.NullableFlag_Empty
	} else {
		// command is not allowed for bundle
		if taskType.IsBundle() {
			return invalidErr(er.ErrFieldNotRequired, f.Name(), f.Value())
		}
		otask.Command.Flag = proto.NullableFlag_Available
		otask.Command.Value = f.Value()
	}

	// validate field 'condition'
	if f, ok := tsk.GetFieldCondition(); !ok {
		otask.Condition.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.Condition.Flag = proto.NullableFlag_Empty
	} else {
		otask.Condition.Flag = proto.NullableFlag_Available
		otask.Condition.Value = f.Value()
	}

	// validate field 'std_out_log'
	if f, ok := tsk.GetFieldOutLogFile(); !ok {
		otask.OutLogFile.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.OutLogFile.Flag = proto.NullableFlag_Empty
	} else {
		if taskType.IsBundle() {
			return invalidErr(er.ErrFieldNotRequired, f.Name(), f.Value())
		}
		otask.OutLogFile.Flag = proto.NullableFlag_Available
		otask.OutLogFile.Value = f.Value()
	}

	// validate field 'std_err_log'
	if f, ok := tsk.GetFieldErrLogFile(); !ok {
		otask.ErrLogFile.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.ErrLogFile.Flag = proto.NullableFlag_Empty
	} else {
		if taskType.IsBundle() {
			return invalidErr(er.ErrFieldNotRequired, f.Name(), f.Value())
		}
		otask.ErrLogFile.Flag = proto.NullableFlag_Available
		otask.ErrLogFile.Value = f.Value()
	}

	// validate field 'label'
	if f, ok := tsk.GetFieldLabel(); !ok {
		otask.Label.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.Label.Flag = proto.NullableFlag_Empty
	} else {
		otask.Label.Flag = proto.NullableFlag_Available
		otask.Label.Value = f.Value()
	}

	// validate field 'profile'
	if f, ok := tsk.GetFieldProfile(); !ok {
		otask.Profile.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.Profile.Flag = proto.NullableFlag_Empty
	} else {
		if taskType.IsBundle() {
			return invalidErr(er.ErrFieldNotRequired, f.Name(), f.Value())
		}
		otask.Profile.Flag = proto.NullableFlag_Available
		otask.Profile.Value = f.Value()
	}

	// validate field 'priority'
	if f, ok := tsk.GetFieldPriority(); !ok {
		otask.Priority.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.Priority.Flag = proto.NullableFlag_Empty
	} else {
		otask.Priority.Flag = proto.NullableFlag_Available
		otask.Priority.Value = f.Value()
	}

	// validate field 'run_days'
	if f, ok := tsk.GetFieldRunDays(); !ok {
		otask.RunDays.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		otask.RunDays.Flag = proto.NullableFlag_Empty
	} else {
		otask.RunDays.Flag = proto.NullableFlag_Available
		otask.RunDays.Value = int32(f.Value())
	}

	var rmST, rmRW, rmSM bool

	// validate field 'start_times'
	if f, ok := tsk.GetFieldStartTimes(); !ok {
		otask.RunDays.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		rmST = true
		otask.StartTimes.Flag = proto.NullableFlag_Empty
	} else {
		otask.StartTimes.Flag = proto.NullableFlag_Available
		otask.StartTimes.Value = f.Value()
	}

	// validate field 'start_mins'
	if f, ok := tsk.GetFieldStartMins(); !ok {
		otask.StartMins.Flag = proto.NullableFlag_NotAvailable
	} else if f.Empty() {
		rmSM = true
		otask.StartMins.Flag = proto.NullableFlag_Empty
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
	} else if f.Empty() {
		rmRW = true
		otask.RunWindow.Flag = proto.NullableFlag_Empty
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

	// check if job is scheduling to window run and if so,
	// it should have both "run_window" and "start_mins" attributes
	if (hasRW != hasSM) && (rmRW != rmSM) {
		return badErr(er.ErrIncompleteWindowRun)
	}

	schedToBatchRun := hasST && !rmST
	schedToWindowRun := hasRW && hasSM && !rmRW && !rmSM

	// check if job is scheduling to both window run and batch run
	if schedToBatchRun && schedToWindowRun {
		return badErr(er.ErrBatchWindowRun)
	}

	// cannot remove start_times for existing job which doesn't have start_times
	if runFlag.IsWindow() && rmST {
		return badErr(er.ErrNonEmptyValueRequired)
	}

	// trying to switch from batch run to window run
	if runFlag.IsBatch() && schedToWindowRun {
		if !rmST {
			return badErr(er.ErrRemoveBatchRun)
		}
	}

	// trying to switch from window run to batch run
	if runFlag.IsWindow() && schedToBatchRun {
		if !rmSM || !rmRW {
			return badErr(er.ErrRemoveWindowRun)
		}
	}

	return nil
}

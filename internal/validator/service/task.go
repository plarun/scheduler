package service

import (
	"errors"

	"github.com/plarun/scheduler/api/types/entity/task"
	er "github.com/plarun/scheduler/internal/validator/errors"
	proto "github.com/plarun/scheduler/proto"
)

type Actioner interface {
	// validate fields

}

func newValidatedTask() *proto.ValidatedTaskEntity {
	return &proto.ValidatedTaskEntity{
		Command:    &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		Condition:  &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		ErrLogFile: &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		Label:      &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		Machine:    &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		OutLogFile: &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		Parent:     &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		Priority:   &proto.NullableInt32{Flag: proto.NullableFlag_NotAvailable},
		Profile:    &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
		RunDays:    &proto.NullableInt32{Flag: proto.NullableFlag_NotAvailable},
		RunWindow:  &proto.NullableTimeRange{Flag: proto.NullableFlag_NotAvailable, Value: &proto.TimeRange{}},
		StartMins:  &proto.NullableInt32S{Flag: proto.NullableFlag_NotAvailable},
		StartTimes: &proto.NullableStrings{Flag: proto.NullableFlag_NotAvailable},
		Type:       &proto.NullableString{Flag: proto.NullableFlag_NotAvailable},
	}
}

func prepareTask(itask *proto.ParsedTaskEntity) (*task.TaskEntity, error) {
	tsk := task.NewTaskEntity(itask.Target)

	// field: command
	if f, ok := itask.Fields[string(task.FIELD_COMMAND)]; ok {
		empty := false
		if len(f) == 0 {
			empty = true
		}
		tsk.SetFieldCommand(f, empty)
	}

	// field: condition
	if f, ok := itask.Fields[string(task.FIELD_CONDITION)]; ok {
		var (
			val   string
			empty bool
			err   error
		)
		if len(f) == 0 {
			empty = true
		} else {
			if val, err = castFieldStartCondition(f); err != nil {
				return nil, err
			}
		}
		tsk.SetFieldCondition(val, empty)
	}

	// field: label
	if f, ok := itask.Fields[string(task.FIELD_LABEL)]; ok {
		empty := false
		if len(f) == 0 {
			empty = true
		}
		tsk.SetFieldLabel(f, empty)
	}

	// field: out_log_file
	if f, ok := itask.Fields[string(task.FIELD_OUT_LOG_FILE)]; ok {
		empty := false
		if len(f) == 0 {
			empty = true
		}
		tsk.SetFieldOutLogFile(f, empty)
	}

	// field: err_log_file
	if f, ok := itask.Fields[string(task.FIELD_ERR_LOG_FILE)]; ok {
		empty := false
		if len(f) == 0 {
			empty = true
		}
		tsk.SetFieldErrLogFile(f, empty)
	}

	// field: machine
	if f, ok := itask.Fields[string(task.FIELD_MACHINE)]; ok {
		empty := false
		if len(f) == 0 {
			empty = true
		}
		tsk.SetFieldMachine(f, empty)
	}

	// field: parent
	if f, ok := itask.Fields[string(task.FIELD_PARENT)]; ok {
		empty := false
		if len(f) == 0 {
			empty = true
		}
		tsk.SetFieldParent(f, empty)
	}

	// field: priority
	if f, ok := itask.Fields[string(task.FIELD_PRIORITY)]; ok {
		var (
			val   int32
			empty bool
			err   error
		)
		if len(f) == 0 {
			empty = true
		} else {
			if val, err = castFieldPriority(f); err != nil {
				return nil, err
			}
		}
		tsk.SetFieldPriority(val, empty)
	}

	// field: profile
	if f, ok := itask.Fields[string(task.FIELD_PROFILE)]; ok {
		empty := false
		if len(f) == 0 {
			empty = true
		}
		tsk.SetFieldProfile(f, empty)
	}

	// field: run_days
	if f, ok := itask.Fields[string(task.FIELD_RUN_DAYS)]; ok {
		var (
			val   int32
			empty bool
			err   error
		)
		if len(f) == 0 {
			empty = true
		} else {
			if val, err = castFieldRundays(f); err != nil {
				return nil, err
			}
		}
		tsk.SetFieldRunDays(val, empty)
	}

	// field: start_times
	if f, ok := itask.Fields[string(task.FIELD_START_TIMES)]; ok {
		val, err := castFieldStartTimes(f)
		empty := false
		if errors.Is(err, er.ErrNonEmptyValueRequired) {
			empty = true
		} else if err != nil {
			return nil, err
		}
		tsk.SetFieldStartTimes(val, empty)
	}

	// field: start_mins
	if f, ok := itask.Fields[string(task.FIELD_START_MINS)]; ok {
		val, err := castFieldStartMins(f)
		empty := false
		if errors.Is(err, er.ErrNonEmptyValueRequired) {
			empty = true
		} else if err != nil {
			return nil, err
		}
		tsk.SetFieldStartMins(val, empty)
	}

	// field: run_window
	if f, ok := itask.Fields[string(task.FIELD_RUN_WINDOW)]; ok {
		start, end, err := castFieldRunWindow(f)
		empty := false
		if errors.Is(err, er.ErrNonEmptyValueRequired) {
			empty = true
		} else if err != nil {
			return nil, err
		}
		tsk.SetFieldRunWindow(start, end, empty)
	}

	// field: type
	if f, ok := itask.Fields[string(task.FIELD_TYPE)]; ok {
		val, err := castFieldType(f)
		empty := false
		if errors.Is(err, er.ErrNonEmptyValueRequired) {
			empty = true
		} else if err != nil {
			return nil, err
		}
		tsk.SetFieldType(val, empty)
	}

	return tsk, nil
}

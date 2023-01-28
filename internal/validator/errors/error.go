package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNonEmptyValueRequired error = errors.New("value is empty")
	ErrFieldRequired         error = errors.New("field is required")
	ErrFieldNotRequired      error = errors.New("field is not required")

	ErrTaskNotExist     error = errors.New("task not exist")
	ErrTaskAlreadyExist error = errors.New("task already exist")
	ErrTaskMaxLength    error = errors.New("task name can only have 64 characters")
	ErrTaskInvalidChar  error = errors.New("task name can only contain alpha numeric and underscore chars")

	ErrBadDeleteTask error = errors.New("delete task should not have other attributes")

	ErrInvalidActionAttr     error = errors.New("action attribute should be either insert/delete/update")
	ErrInvalidTypeAttr       error = errors.New("type attribute should be either bundle/callable")
	ErrInvalidPriorityAttr   error = errors.New("priority should be either (low=0, normal=1, important=2, critical=3)")
	ErrInvalidRundaysAttr    error = errors.New("run_days attribute should be comma separated any values of su,mo,tu,we,th,fr,sa")
	ErrInvalidStartTimesAttr error = errors.New("start_times attribute should be comma separated values of time in format hh:mm or hh:mm:ss")
	ErrInvalidRunWindowAttr  error = errors.New("run_window attribute should only have start and end time separated by '-' and in format hh:mm or hh:mm:ss")
	ErrInvalidStartMinsAttr  error = errors.New("start_min attribute should be comma separated values of minutes")

	ErrRepeatedRundaysAttr error = errors.New("run_days attribute contains repeated days")

	ErrBatchWindowRun      error = errors.New("task should not have both batch run and window run")
	ErrIncompleteWindowRun error = errors.New("window run should have both 'run_window' and 'start_mins' attributes")
	ErrRemoveBatchRun      error = errors.New("start_times should be set to null to have window run")
	ErrRemoveWindowRun     error = errors.New("run_window and start_mins should be set to null to have batch run")

	ErrNonEditableTypeAttr error = errors.New("type attribute cannot be changed")
)

type InvalidTaskFieldError struct {
	Action string
	Target string
	Field  string
	Value  string
	Err    error
}

func (e *InvalidTaskFieldError) Error() string {
	return fmt.Sprintf("action=\"%s\" target=\"%s\" field=\"%s\" value=\"%s\" error=\"%s\"",
		e.Action, e.Target, e.Field, e.Value, e.Err.Error())
}

type BadTaskDefError struct {
	Action string
	Target string
	Err    error
}

func (e *BadTaskDefError) Error() string {
	return fmt.Sprintf("action=\"%s\" target=\"%s\" error=\"%s\"",
		e.Action, e.Target, e.Err.Error())
}

type InternalSqlError struct {
	Query string
	Err   error
}

func (e *InternalSqlError) Error() string {
	return fmt.Sprintf("query=\"%s\" error=\"%s\"", e.Query, e.Err)
}

package error

import (
	"errors"
	"fmt"
)

var (
	ErrTaskNotFound error = errors.New("task not found")

	ErrEventGreenOnUnstableTask  error = errors.New("event green is not effective on active task")
	ErrEventRedOnUnstableTask    error = errors.New("event red is not effective on active task")
	ErrEventFreezeOnUnstableTask error = errors.New("event freeze is not effective on active task")
	ErrEventResetOnUnstableTask  error = errors.New("event reset is not effective on active task")
	ErrEventStartOnUnstableTask  error = errors.New("event start is not effective on active task")
)

type TaskNotFoundError struct {
	Id   int64
	Name string
	isId bool
}

func NewTaskNotFoundForIdError(id int64) *TaskNotFoundError {
	return &TaskNotFoundError{
		Id:   id,
		isId: true,
	}
}

func NewTaskNotFoundForNameError(name string) *TaskNotFoundError {
	return &TaskNotFoundError{
		Name: name,
		isId: false,
	}
}

func (e TaskNotFoundError) Error() string {
	if e.isId {
		return fmt.Sprintf("Task for id - %d not found", e.Id)
	}
	return fmt.Sprintf("Task - %s not found", e.Name)
}

type TaskEventError struct {
	Msg error
}

func NewTaskEventError(msg error) *TaskEventError {
	return &TaskEventError{
		Msg: msg,
	}
}

func (e TaskEventError) Error() string {
	return e.Msg.Error()
}

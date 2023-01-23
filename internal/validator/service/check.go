package service

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/validator/errors"
	"github.com/plarun/scheduler/proto"
)

type TaskValidater interface {
	Validate() error
	Get() *proto.ValidatedTaskEntity
}

func ValidateTaskAction(tsk *proto.ParsedTaskEntity) (*proto.ValidatedTaskEntity, error) {
	action, err := castFieldAction(tsk.Action)
	if err != nil {
		return nil, fmt.Errorf("ValidateJob: %w", err)
	}

	var vld TaskValidater

	switch task.Action(action) {
	case task.ActionInsert:
		vld = newInsertValidation(tsk)
	case task.ActionUpdate:
		vld = newUpdateValidation(tsk)
	case task.ActionDelete:
		vld = newDeleteValidation(tsk)
	default:
		return nil, errors.ErrInvalidActionAttr
	}

	if err := vld.Validate(); err != nil {
		return nil, fmt.Errorf("ValidateJob: %w", err)
	}

	return vld.Get(), nil
}

package handler

import (
	"errors"

	es "github.com/plarun/scheduler/internal/eventserver/service"
	"github.com/plarun/scheduler/proto"
	"golang.org/x/net/context"
)

// TaskExecService is a grpc server for communicating with worker service
type TaskService struct {
	proto.UnimplementedTaskServiceServer
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (s TaskService) GetDefinition(ctx context.Context, req *proto.TaskDefinitionRequest) (*proto.TaskDefinitionResponse, error) {
	res := &proto.TaskDefinitionResponse{
		IsValid: false,
		Task:    &proto.TaskDefinition{},
	}

	// process the validated task entities' actions
	if tasks, err := es.GetTaskDefinition(ctx, req.TaskName); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		return res, err
	} else {
		res.IsValid = true
		res.Task = tasks
		return res, nil
	}
}

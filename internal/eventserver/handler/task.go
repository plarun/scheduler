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

	// get task definition
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

func (s TaskService) GetStatus(ctx context.Context, req *proto.TaskLatestStatusRequest) (*proto.TaskLatestStatusResponse, error) {
	res := &proto.TaskLatestStatusResponse{
		IsValid: false,
		Status:  &proto.TaskRunStatus{},
	}

	// get task latest status
	if st, err := es.GetTaskLatestStatus(ctx, req.TaskName); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		return res, err
	} else {
		res.IsValid = true
		res.Status = st
		return res, nil
	}
}

func (s TaskService) GetRuns(ctx context.Context, req *proto.TaskRunsRequest) (*proto.TaskRunsResponse, error) {
	res := &proto.TaskRunsResponse{
		IsValid: false,
		Runs:    make([]*proto.TaskRunStatus, 0),
	}

	// get task latest status
	if runs, err := es.GetTaskRuns(ctx, req.TaskName, req.Count, req.RunDate); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		return res, err
	} else {
		res.IsValid = true
		res.Runs = runs
		return res, nil
	}
}

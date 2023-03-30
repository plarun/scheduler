package handler

import (
	"errors"

	es "github.com/plarun/scheduler/internal/eventserver/service"
	er "github.com/plarun/scheduler/pkg/error"
	"github.com/plarun/scheduler/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	tasks, err := es.GetTaskDefinition(ctx, req.TaskName)
	if err != nil {
		if errors.Is(err, er.DatabaseError{}) {
			return res, status.Errorf(codes.Internal, "Internal Server error")
		}
		return res, status.Errorf(codes.Unknown, "Unknown error")
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
	st, err := es.GetTaskLatestStatus(ctx, req.TaskName)
	if err != nil {
		if errors.Is(err, er.DatabaseError{}) {
			return res, status.Errorf(codes.Internal, "Internal Server error")
		}
		return res, status.Errorf(codes.Unknown, "Unknown error")
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
	runs, err := es.GetTaskRuns(ctx, req.TaskName, req.Count, req.RunDate)
	if err != nil {
		if errors.Is(err, er.DatabaseError{}) {
			return res, status.Errorf(codes.Internal, "Internal Server error")
		}
		return res, status.Errorf(codes.Unknown, "Unknown error")
	} else {
		res.IsValid = true
		res.Runs = runs
		return res, nil
	}
}

func (s TaskService) SendEvent(ctx context.Context, req *proto.TaskEventRequest) (*proto.TaskEventResponse, error) {
	res := &proto.TaskEventResponse{}

	// get task latest status
	r, err := es.ActionTaskEvent(ctx, req.TaskName, req.Event)
	if err != nil {
		if errors.Is(err, er.DatabaseError{}) {
			return res, status.Errorf(codes.Internal, "Internal Server error")
		}
		return res, status.Errorf(codes.Unknown, "Unknown error")
	}
	if !r.IsSuccess {
		res.Success = false
		res.Msg = r.Msg.Error()
		return res, nil
	} else {
		res.Success = true
		res.Msg = ""
		return res, nil
	}
}

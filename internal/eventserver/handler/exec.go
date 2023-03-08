package handler

import (
	"errors"

	"github.com/plarun/scheduler/api/types/entity/task"
	es "github.com/plarun/scheduler/internal/eventserver/service"
	"github.com/plarun/scheduler/proto"
	"golang.org/x/net/context"
)

// TaskExecService is a grpc server for communicating with worker service
type TaskExecService struct {
	proto.UnimplementedTaskExecServiceServer
}

func NewTaskExecService() *TaskExecService {
	return &TaskExecService{}
}

func (s TaskExecService) PullReadyTasks(ctx context.Context, req *proto.ReadyTasksPullRequest) (*proto.ReadyTasksPullResponse, error) {
	res := &proto.ReadyTasksPullResponse{
		TaskIds: make([]int64, 0),
	}

	// process the validated task entities' actions
	if tasks, err := es.PullReadyTasks(ctx); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}

		return res, err
	} else {
		res.TaskIds = tasks
		return res, nil
	}
}

func (s TaskExecService) GetTask(ctx context.Context, req *proto.TaskInfoRequest) (*proto.TaskInfoResponse, error) {
	res := &proto.TaskInfoResponse{
		TaskId:     req.TaskId,
		Command:    "",
		OutLogFile: "",
		ErrLogFile: "",
	}

	// process the validated task entities' actions
	if cmd, fout, ferr, err := es.GetTaskCommand(ctx, req.TaskId); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		return res, err
	} else {
		res.Command = cmd
		res.OutLogFile = fout
		res.ErrLogFile = ferr
		return res, nil
	}
}

func (s TaskExecService) SetTaskStatus(ctx context.Context, req *proto.TaskStatusRequest) (*proto.TaskStatusResponse, error) {
	res := &proto.TaskStatusResponse{}

	var state task.State
	switch req.Status {
	case proto.TaskStatus_RUNNING:
		state = task.StateRunning
	case proto.TaskStatus_SUCCESS:
		state = task.StateSuccess
	case proto.TaskStatus_FAILURE:
		state = task.StateFailure
	case proto.TaskStatus_ABORTED:
		state = task.StateAborted
	}

	// process the validated task entities' actions
	if err := es.ChangeTaskState(ctx, req.TaskId, state); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		return res, err
	} else {
		return res, nil
	}
}

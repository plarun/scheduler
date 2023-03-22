package handler

import (
	"errors"

	es "github.com/plarun/scheduler/internal/eventserver/service"
	"github.com/plarun/scheduler/proto"
	"golang.org/x/net/context"
)

// WaitingTaskAwakeService is a grpc server for communicating with eventserver service
type WaitingTaskAwakeService struct {
	proto.UnimplementedWaitTaskServiceServer
}

func NewTaskExecService() *WaitingTaskAwakeService {
	return &WaitingTaskAwakeService{}
}

func (w WaitingTaskAwakeService) Awake(ctx context.Context, req *proto.DependentTaskAwakeRequest) (*proto.DependentTaskAwakeResponse, error) {
	res := &proto.DependentTaskAwakeResponse{}

	// process the validated task entities' actions
	if err := es.AwakeWaitingDependentTasks(ctx, req.TaskId); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		return res, err
	} else {
		return res, nil
	}
}

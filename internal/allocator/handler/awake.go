package handler

import (
	"errors"
	"log"

	alloc "github.com/plarun/scheduler/internal/allocator/service"
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

	log.Printf("Request to Awake the task id - %d", req.TaskId)

	// process the validated task entities' actions
	if err := alloc.AwakeWaitingDependentTasks(ctx, req.TaskId); err != nil {
		for errors.Unwrap(err) != nil {
			err = errors.Unwrap(err)
		}
		return res, err
	} else {
		return res, nil
	}
}

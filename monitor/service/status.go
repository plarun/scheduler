package service

import (
	"context"

	pb "github.com/plarun/scheduler/monitor/data"
	"github.com/plarun/scheduler/monitor/locker"
	"google.golang.org/grpc"
)

var statusService *StatusService = nil

type StatusService struct {
	pb.UnimplementedUpdateStatusServer
	eventServerClient pb.UpdateStatusClient
}

func InitUpdateStatusClient(conn *grpc.ClientConn) {
	statusService = &StatusService{
		eventServerClient: pb.NewUpdateStatusClient(conn),
	}
}

func GetStatusService() *StatusService {
	return statusService
}

func (stat StatusService) Update(ctx context.Context, req *pb.UpdateStatusReq) (*pb.UpdateStatusRes, error) {
	if _, err := stat.eventServerClient.Update(ctx, req); err != nil {
		return &pb.UpdateStatusRes{}, err
	}

	updatedStatus := req.GetStatus()
	jobName := req.GetJobName()

	if updatedStatus == pb.NewStatus_CHANGE_READY {
		locker.GetLocker().Put(jobName)
	} else if updatedStatus == pb.NewStatus_CHANGE_FAILED || updatedStatus == pb.NewStatus_CHANGE_SUCCESS {
		locker.GetLocker().Free(jobName)
	}

	return &pb.UpdateStatusRes{}, nil
}
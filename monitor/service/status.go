package service

import (
	"context"

	pb "github.com/plarun/scheduler/monitor/data"
	"github.com/plarun/scheduler/monitor/locker"
	"google.golang.org/grpc"
)

// Singleton instance of StatusService
var statusService *StatusService = nil

type StatusService struct {
	pb.UnimplementedUpdateStatusServer
	eventServerClient pb.UpdateStatusClient
}

// InitUpdateStatusClient initiates the StatusService
func InitUpdateStatusClient(conn *grpc.ClientConn) {
	statusService = &StatusService{
		eventServerClient: pb.NewUpdateStatusClient(conn),
	}
}

// GetStatusService returns singleton instance of StatusService
func GetStatusService() *StatusService {
	return statusService
}

// Update handles the job status update
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

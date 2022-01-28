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
	client pb.UpdateStatusClient
}

func InitUpdateStatusClient(conn *grpc.ClientConn) {
	statusService = &StatusService{
		client: pb.NewUpdateStatusClient(conn),
	}
}

func GetStatusService() *StatusService {
	return statusService
}

func (stat StatusService) Update(ctx context.Context, req *pb.UpdateStatusReq) (*pb.UpdateStatusRes, error) {
	_, err := stat.client.Update(ctx, req)
	if err != nil {
		return &pb.UpdateStatusRes{}, err
	}

	updatedStatus := req.GetStatus()
	if updatedStatus == pb.NewStatus_CHANGE_READY {
		locker := locker.GetLocker()
		locker.Put(req.GetJobName())
	}
	return &pb.UpdateStatusRes{}, nil
}

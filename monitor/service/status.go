package service

import (
	"context"

	pb "github.com/plarun/scheduler/monitor/data"
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
	return &pb.UpdateStatusRes{}, nil
}

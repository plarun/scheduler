package service

import (
	"context"

	pb "github.com/plarun/scheduler/monitor/data"
	"google.golang.org/grpc"
)

// Singleton instance of ConditionService
var conditionService *ConditionService = nil

type ConditionService struct {
	pb.UnimplementedConditionServer
	pickerClient pb.ConditionClient
}

// InitConditionClient initiates the ConditionService
func InitConditionClient(pickerConn *grpc.ClientConn) {
	conditionService = &ConditionService{
		pickerClient: pb.NewConditionClient(pickerConn),
	}
}

// GetConditionService returns singleton instance of ConditionService
func GetConditionService() *ConditionService {
	return conditionService
}

// ConditionStatus passes the successors status to the picker to release the waiting jobs
func (cond ConditionService) ConditionStatus(ctx context.Context, req *pb.JobConditionReq) (*pb.JobConditionRes, error) {
	if _, err := cond.pickerClient.ConditionStatus(ctx, req); err != nil {
		return &pb.JobConditionRes{}, err
	}

	return &pb.JobConditionRes{}, nil
}

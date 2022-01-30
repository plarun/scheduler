package service

import (
	"context"
	"log"

	pb "github.com/plarun/scheduler/monitor/data"
	"google.golang.org/grpc"
)

var conditionService *ConditionService = nil

type ConditionService struct {
	pb.UnimplementedConditionServer
	pickerClient pb.ConditionClient
}

func InitConditionClient(pickerConn *grpc.ClientConn) {
	conditionService = &ConditionService{
		pickerClient: pb.NewConditionClient(pickerConn),
	}
}

func GetConditionService() *ConditionService {
	return conditionService
}

func (cond ConditionService) ConditionStatus(ctx context.Context, req *pb.JobConditionReq) (*pb.JobConditionRes, error) {
	log.Println("ConditionService.ConditionStatus")
	_, err := cond.pickerClient.ConditionStatus(ctx, req)
	if err != nil {
		return &pb.JobConditionRes{}, err
	}

	return &pb.JobConditionRes{}, nil
}

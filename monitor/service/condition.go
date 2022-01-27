package service

import (
	pb "github.com/plarun/scheduler/monitor/data"
	"google.golang.org/grpc"
)

var conditionService *ConditionService = nil

type ConditionService struct {
	pb.UnimplementedConditionServer
	eventServerClient pb.ConditionClient
	pickerClient      pb.ConditionClient
}

func InitConditionClient(eventServerConn *grpc.ClientConn, pickerConn *grpc.ClientConn) {
	conditionService = &ConditionService{
		eventServerClient: pb.NewConditionClient(eventServerConn),
		pickerClient:      pb.NewConditionClient(pickerConn),
	}
}

func GetConditionService() *ConditionService {
	return conditionService
}

package client

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type statusGrpcConnection struct {
	addr    string
	client  proto.TaskServiceClient
	conn    *grpc.ClientConn
	request *proto.TaskLatestStatusRequest
}

func NewStatusGrpcConnection(addr string, req *proto.TaskLatestStatusRequest) grpcconn.GrpcConnecter {
	return &statusGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (s *statusGrpcConnection) Connect() error {
	conn, err := grpc.Dial(s.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	s.conn = conn
	s.client = proto.NewTaskServiceClient(conn)
	return nil
}

func (s *statusGrpcConnection) Request() (interface{}, error) {
	if s.client == nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := s.client.GetStatus(context.Background(), s.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *statusGrpcConnection) Close() error {
	if s.client != nil {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

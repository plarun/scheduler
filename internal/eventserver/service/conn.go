package service

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type waitTaskGrpcConnection struct {
	addr    string
	client  proto.WaitTaskServiceClient
	conn    *grpc.ClientConn
	request *proto.DependentTaskAwakeRequest
}

func NewWaitTaskGrpcConnection(addr string, req *proto.DependentTaskAwakeRequest) grpcconn.GrpcConnecter {
	return &waitTaskGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (t *waitTaskGrpcConnection) Connect() error {
	conn, err := grpc.Dial(t.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	t.conn = conn
	t.client = proto.NewWaitTaskServiceClient(conn)
	return nil
}

func (t *waitTaskGrpcConnection) Request() (interface{}, error) {
	if t.client == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := t.client.Awake(context.Background(), t.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *waitTaskGrpcConnection) Close() error {
	if t.client != nil {
		if err := t.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

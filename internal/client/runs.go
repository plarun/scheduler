package client

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type taskRunsGrpcConnection struct {
	addr    string
	client  proto.TaskServiceClient
	conn    *grpc.ClientConn
	request *proto.TaskRunsRequest
}

func NewTaskRunsGrpcConnection(addr string, req *proto.TaskRunsRequest) grpcconn.GrpcConnecter {
	return &taskRunsGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (t *taskRunsGrpcConnection) Connect() error {
	conn, err := grpc.Dial(t.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	t.conn = conn
	t.client = proto.NewTaskServiceClient(conn)
	return nil
}

func (t *taskRunsGrpcConnection) Request() (interface{}, error) {
	if t.client == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := t.client.GetRuns(context.Background(), t.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *taskRunsGrpcConnection) Close() error {
	if t.client != nil {
		if err := t.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

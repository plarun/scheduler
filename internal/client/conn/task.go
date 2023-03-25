package conn

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type taskGrpcConnection struct {
	addr    string
	client  proto.TaskServiceClient
	conn    *grpc.ClientConn
	request *proto.TaskDefinitionRequest
}

func NewTaskGrpcConnection(addr string, req *proto.TaskDefinitionRequest) grpcconn.GrpcConnecter {
	return &taskGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (t *taskGrpcConnection) Connect() error {
	conn, err := grpc.Dial(t.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	t.conn = conn
	t.client = proto.NewTaskServiceClient(conn)
	return nil
}

func (t *taskGrpcConnection) Request() (interface{}, error) {
	if t.client == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := t.client.GetDefinition(context.Background(), t.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *taskGrpcConnection) Close() error {
	if t.client != nil {
		if err := t.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

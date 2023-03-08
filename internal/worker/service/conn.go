package service

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type taskInfoGrpcConnection struct {
	addr    string
	client  proto.TaskExecServiceClient
	conn    *grpc.ClientConn
	request *proto.TaskInfoRequest
}

func NewTaskInfoGrpcConnection(addr string, req *proto.TaskInfoRequest) grpcconn.GrpcConnecter {
	return &taskInfoGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (t *taskInfoGrpcConnection) Connect() error {
	conn, err := grpc.Dial(t.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	t.conn = conn
	t.client = proto.NewTaskExecServiceClient(conn)
	return nil
}

func (t *taskInfoGrpcConnection) Request() (interface{}, error) {
	if t.client == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := t.client.GetTask(context.Background(), t.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *taskInfoGrpcConnection) Close() error {
	if t.client != nil {
		if err := t.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

type taskStatusGrpcConnection struct {
	addr    string
	client  proto.TaskExecServiceClient
	conn    *grpc.ClientConn
	request *proto.TaskStatusRequest
}

func NewTaskStatusGrpcConnection(addr string, req *proto.TaskStatusRequest) grpcconn.GrpcConnecter {
	return &taskStatusGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (t *taskStatusGrpcConnection) Connect() error {
	conn, err := grpc.Dial(t.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	t.conn = conn
	t.client = proto.NewTaskExecServiceClient(conn)
	return nil
}

func (t *taskStatusGrpcConnection) Request() (interface{}, error) {
	if t.client == nil {
		if err := t.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := t.client.SetTaskStatus(context.Background(), t.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (t *taskStatusGrpcConnection) Close() error {
	if t.client != nil {
		if err := t.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

type readyTasksGrpcConnection struct {
	addr    string
	client  proto.TaskExecServiceClient
	conn    *grpc.ClientConn
	request *proto.ReadyTasksPullRequest
}

func NewReadyTasksGrpcConnection(addr string, req *proto.ReadyTasksPullRequest) grpcconn.GrpcConnecter {
	return &readyTasksGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (r *readyTasksGrpcConnection) Connect() error {
	conn, err := grpc.Dial(r.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	r.conn = conn
	r.client = proto.NewTaskExecServiceClient(conn)
	return nil
}

func (r *readyTasksGrpcConnection) Request() (interface{}, error) {
	if r.client == nil {
		if err := r.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := r.client.PullReadyTasks(context.Background(), r.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *readyTasksGrpcConnection) Close() error {
	if r.client != nil {
		if err := r.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

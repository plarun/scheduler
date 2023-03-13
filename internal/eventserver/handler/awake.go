package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type dependentAwakeGrpcConnection struct {
	addr    string
	client  proto.WaitTaskServiceClient
	conn    *grpc.ClientConn
	request *proto.DependentTaskAwakeRequest
}

func NewDependentAwakeGrpcConnection(addr string, req *proto.DependentTaskAwakeRequest) grpcconn.GrpcConnecter {
	return &dependentAwakeGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (d *dependentAwakeGrpcConnection) Connect() error {
	conn, err := grpc.Dial(d.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Connect: failed to get grpc client connnection")
	}
	d.conn = conn
	d.client = proto.NewWaitTaskServiceClient(conn)
	return nil
}

func (d *dependentAwakeGrpcConnection) Request() (interface{}, error) {
	if d.client == nil {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := d.client.Awake(context.Background(), d.request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d *dependentAwakeGrpcConnection) Close() error {
	if d.client != nil {
		if err := d.conn.Close(); err != nil {
			return fmt.Errorf("Close: failed to close grpc client connection")
		}
	}
	return nil
}

func awakeWaitingDependentTasks(ctx context.Context, id int64, state task.State) error {
	// invoke the dependent tasks if any of them are waiting
	if state.IsSuccess() || state.IsFailure() {
		// get all the dependent task ids of this task
		// pass them to check on sched_wait
		// then awake them if any

		// rpc to allocator to awake dep tasks from waiting
		req := &proto.DependentTaskAwakeRequest{
			TaskId: id,
		}

		addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.Allocator.Port)
		allocConn := NewDependentAwakeGrpcConnection(addr, req)

		log.Println("Routing the request to awake waiting dependent task from allocator")

		if err := allocConn.Connect(); err != nil {
			return err
		}

		valRes, err := allocConn.Request()
		if err != nil {
			return err
		}

		_, ok := valRes.(*proto.DependentTaskAwakeResponse)
		if !ok {
			return fmt.Errorf("internal err")
		}

		if err := allocConn.Close(); err != nil {
			return err
		}
		return nil
	}
	return nil
}
package handler

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type validationGrpcConnection struct {
	addr    string
	client  proto.ValidatedActionServiceClient
	conn    *grpc.ClientConn
	request *proto.ParsedEntitiesRequest
}

func NewValidationGrpcConnection(addr string, req *proto.ParsedEntitiesRequest) grpcconn.GrpcConnecter {
	return &validationGrpcConnection{
		addr:    addr,
		request: req,
	}
}

func (v *validationGrpcConnection) Connect() error {
	conn, err := grpc.Dial(v.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Connect: failed to get grpc client connnection")
	}
	v.conn = conn
	v.client = proto.NewValidatedActionServiceClient(conn)
	return nil
}

func (v *validationGrpcConnection) Request() (interface{}, error) {
	if v.client == nil {
		if err := v.Connect(); err != nil {
			return nil, err
		}
	}

	res, err := v.client.Route(context.Background(), v.request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (v *validationGrpcConnection) Close() error {
	if v.client != nil {
		if err := v.conn.Close(); err != nil {
			return fmt.Errorf("Close: failed to close grpc client connection")
		}
	}
	return nil
}

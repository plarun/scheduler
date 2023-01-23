package client

import (
	"context"
	"fmt"

	"github.com/plarun/scheduler/pkg/grpcconn"
	"github.com/plarun/scheduler/proto"
	"google.golang.org/grpc"
)

type definitionGrpcConnection struct {
	addr      string
	client    proto.ParsedActionServiceClient
	conn      *grpc.ClientConn
	onlyCheck bool
	request   *proto.ParsedEntitiesRequest
}

func NewDefinitionGrpcConnection(addr string, flag bool, req *proto.ParsedEntitiesRequest) grpcconn.GrpcConnecter {
	return &definitionGrpcConnection{
		addr:      addr,
		onlyCheck: flag,
		request:   req,
	}
}

func (d *definitionGrpcConnection) Connect() error {
	conn, err := grpc.Dial(d.addr, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to get grpc client connnection")
	}
	d.conn = conn
	d.client = proto.NewParsedActionServiceClient(conn)
	return nil
}

func (d *definitionGrpcConnection) Request() (interface{}, error) {
	if d.client == nil {
		if err := d.Connect(); err != nil {
			return nil, err
		}
	}

	d.request.OnlyValidate = d.onlyCheck

	res, err := d.client.Submit(context.Background(), d.request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (d *definitionGrpcConnection) Close() error {
	if d.client != nil {
		if err := d.conn.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection")
		}
	}
	return nil
}

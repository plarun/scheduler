package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/plarun/scheduler/monitor/data"
	"github.com/plarun/scheduler/monitor/locker"
	"github.com/plarun/scheduler/monitor/service"
	"google.golang.org/grpc"
)

const port = 5558

func main() {
	// client connection to Event Server
	eventServerConn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer eventServerConn.Close()

	// client connection to Picker
	pickerConn, err := grpc.Dial("localhost:5556", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer pickerConn.Close()

	service.InitUpdateStatusClient(eventServerConn)
	service.InitConditionClient(pickerConn)

	serve()
}

func serve() {
	addr := fmt.Sprintf(":%d", port)
	// Server listens on tcp port
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	// grpc server can have arguments for unary and stream as server options
	grpcServer := grpc.NewServer()

	// register all servers here
	pb.RegisterUpdateStatusServer(grpcServer, service.GetStatusService())
	pb.RegisterConditionServer(grpcServer, service.GetConditionService())
	pb.RegisterCheckLockServer(grpcServer, locker.NewCheckLockServer())

	log.Printf("Monitor grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

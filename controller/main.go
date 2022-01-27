package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/plarun/scheduler/controller/data"
	"github.com/plarun/scheduler/controller/executor"
	"github.com/plarun/scheduler/controller/queue"
	"github.com/plarun/scheduler/controller/receiver"
	"google.golang.org/grpc"
)

const (
	port = 5557
)

func main() {

	go func() {
		que := queue.GetProcessQueue()
		executor.GetExecutorPool().Start(que)
	}()

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
	pb.RegisterPassJobsServer(grpcServer, receiver.NewPassJobsServer())

	log.Printf("Controller grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

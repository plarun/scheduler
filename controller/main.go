package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/plarun/scheduler/controller/data"
	"github.com/plarun/scheduler/controller/executor"
	"google.golang.org/grpc"
)

const (
	port = 5557
)

func main() {
	// client connection to Monitor
	conn, err := grpc.Dial("localhost:5558", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	statusClient := pb.NewUpdateStatusClient(conn)
	executor.InitUpdateStatusClient(statusClient)

	go func() {
		if err := executor.GetExecutorPool().Start(); err != nil {
			log.Fatal(err)
		}
	}()

	serve()
}

// serve serves the requests
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
	pb.RegisterPassJobsServer(grpcServer, executor.NewPassJobsServer())

	log.Printf("Controller grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

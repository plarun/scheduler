package main

import (
	"fmt"
	"log"
	"net"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
	"github.com/plarun/scheduler/event-server/service"
	"google.golang.org/grpc"
)

const port = 5555

func main() {
	// Connect to sql database
	query.ConnectDB()

	// event server service
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
	pb.RegisterSubmitJilServer(grpcServer, service.JilServer{Database: query.GetDatabase()})
	pb.RegisterNextJobsServer(grpcServer, service.NextJobsServer{Database: query.GetDatabase()})
	pb.RegisterJobStatusServer(grpcServer, service.StatusServer{Database: query.GetDatabase()})

	log.Printf("Scheduler grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

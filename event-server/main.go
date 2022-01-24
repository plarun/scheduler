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

	database := query.GetDatabase()

	// grpc server can have arguments for unary and stream as server options
	grpcServer := grpc.NewServer()

	// register all servers here
	pb.RegisterSubmitJilServer(grpcServer, service.JilServer{Database: database})
	pb.RegisterPickJobsServer(grpcServer, service.NextJobsServer{Database: database})
	pb.RegisterJobStatusServer(grpcServer, service.StatusServer{Database: database})
	pb.RegisterSendEventServer(grpcServer, service.SendEventServer{Database: database})

	log.Printf("Event-Server grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

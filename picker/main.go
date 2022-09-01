package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/plarun/scheduler/picker/checker"
	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/pickpass"

	"google.golang.org/grpc"
)

const port = 5556

func main() {
	// pick client
	pickConn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}

	defer func(pickConn *grpc.ClientConn) {
		err := pickConn.Close()
		if err != nil {

		}
	}(pickConn)

	pickClient := pb.NewPickJobsClient(pickConn)

	// pass client
	passConn, err := grpc.Dial("localhost:5557", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}

	defer func(passConn *grpc.ClientConn) {
		err := passConn.Close()
		if err != nil {

		}
	}(passConn)

	passClient := pb.NewPassJobsClient(passConn)

	jobPickPass := pickpass.InitPickPasser(pickClient, passClient)

	// pick jobs from event-server
	go func() {
		log.Println("Picker started...")

		for ; true; time.Sleep(time.Second * 2) {
			if err := jobPickPass.PickJobs(); err != nil {
				log.Println("Unable to connect to event-server", "retrying...")
				time.Sleep(time.Second * 5)
			}
		}
	}()

	serve()
}

// serve requests
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
	pb.RegisterConditionServer(grpcServer, checker.NewHoldChecker())

	log.Printf("PickPass grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

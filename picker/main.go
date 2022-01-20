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
	log.Println("Picker started...")

	// server service to communicate with controller
	serve()

	// pick client
	pickConn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer pickConn.Close()
	pickClient := pb.NewPickJobsClient(pickConn)

	// pass client
	passConn, err := grpc.Dial("localhost:5557", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer passConn.Close()
	passClient := pb.NewPassJobsClient(passConn)

	jobPickPass := pickpass.NewJobPicker(pickClient, passClient)

	// pick jobs from event-server
	pickErrChan := make(chan error)
	go func() {
		log.Println("client services starting...")
		for {
			if err := jobPickPass.PickJobs(); err != nil {
				pickErrChan <- err
			}
			jobPickPass.Queue.Print()
			time.Sleep(time.Second * 5)
		}
	}()

	log.Fatal(<-pickErrChan)

	// pass jobs to controller
	passErrChan := make(chan error)
	go func() {

		for {
			if err := jobPickPass.PassJobs(); err != nil {
				passErrChan <- err
			}
			time.Sleep(time.Second * 5)
		}
	}()

	log.Fatal(<-passErrChan)
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
	pb.RegisterConditionServer(grpcServer, checker.NewHoldChecker())

	fmt.Printf("Scheduler grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

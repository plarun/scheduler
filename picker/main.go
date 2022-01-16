package main

import (
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/pass"
	"github.com/plarun/scheduler/picker/picker"

	"google.golang.org/grpc"
)

const port = 5556

func main() {
	log.Println("Picker started...")

	// server service to communicate with controller
	serve()

	// client service to communicate with event-server
	clientErrChan := make(chan error)
	go func() {
		conn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("connection failed: %v", err)
		}
		defer conn.Close()

		client := pb.NewNextJobsClient(conn)
		jobPicker := picker.NewJobPicker(client)
		log.Println("client services starting...")
		for {
			time.Sleep(time.Second * 5)
			if err := jobPicker.NextJobs(); err != nil {
				clientErrChan <- err
			}
			jobPicker.Queue.Print()
		}
	}()

	log.Fatal(<-clientErrChan)
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
	pb.RegisterPassJobsServer(grpcServer, pass.NewJobPass())

	fmt.Printf("Scheduler grpc server is running at port: %d\n", port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to start server %v", err)
	}
}

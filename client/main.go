package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/plarun/scheduler/client/data"
	"github.com/plarun/scheduler/client/model"
	"google.golang.org/grpc"

	"github.com/plarun/scheduler/client/job"
)

func main() {
	startClient()
}

// startClient Starts and creates all required client services.
func startClient() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand")
		os.Exit(1)
	}

	conn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	// client services
	clientServices := model.NewClientServices()
	// clientServices.SubmitJil = pb.NewSubmitJilClient(conn)
	clientServices.SendEvent = pb.NewSendEventClient(conn)
	clientServices.JobStatus = pb.NewJobStatusClient(conn)
	clientServices.Dependent = pb.NewJobDependsClient(conn)

	switch cmd := flag.Arg(0); cmd {
	case "submitjil":
		err = submitJil(pb.NewSubmitJilClient(conn))
	case "sendevent":
		err = sendevent(context.Background(), nil)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// sendevent is a subcommand to send a job action event request
func sendevent(ctx context.Context, client interface{}) error {
	fmt.Println("send event")
	fmt.Println("not implemented")
	return nil
}

// submitjil is a subcommand to parse and request the job definitions
func submitJil(client pb.SubmitJilClient) error {
	if flag.NArg() != 2 {
		return fmt.Errorf("invalid argument\nusage:\n\tsubmitjil filepath")
	}

	inputFilename := flag.Arg(1)

	jobInfo := job.NewJobInfo(client)
	if err := jobInfo.SubmitJil(inputFilename); err != nil {
		return err
	}

	return nil
}

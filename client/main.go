package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/plarun/scheduler/client/data"
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

	cmd := flag.Arg(0)

	if cmd == "submitjil" {
		err = submitJil(pb.NewSubmitJilClient(conn))
	} else if cmd == "sendevent" {
		err = sendEvent(pb.NewSendEventClient(conn))
	} else if cmd == "jobstat" || cmd == "jobdef" || cmd == "jobhist" {
		err = status(cmd, pb.NewJobStatusClient(conn))
	} else if cmd == "jobdepend" || cmd == "jobfuture" {
		err = dependents(cmd, pb.NewJobDependsClient(conn))
	} else {
		printHelp()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// sendevent is a subcommand to send a job action event request
func sendEvent(client interface{}) error {
	if flag.NArg() != 3 {
		return fmt.Errorf("invalid argument\nusage:\n\tsendevent <job_name> <event_type>")
	}

	jobName := flag.Arg(1)
	eventType := flag.Arg(2)
	return nil
}

// submitjil is a subcommand to parse and request the job definitions
func submitJil(client pb.SubmitJilClient) error {
	if flag.NArg() != 2 {
		return fmt.Errorf("invalid argument\nusage:\n\tsubmitjil <file_path>")
	}

	inputFilename := flag.Arg(1)

	controller := job.NewJobInfoController(client)
	if err := controller.SubmitJil(inputFilename); err != nil {
		return err
	}

	return nil
}

// status is a subcommand to view latest runs, jobs definition and run history
func status(subCommand string, client pb.JobStatusClient) error {
	controller := job.NewJobStatusController(client)

	if subCommand == "jobstat" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\tjobstat <job_name>")
		}

		jobName := flag.Arg(1)
		err := controller.PrintJobStatus(jobName)
		if err != nil {
			return err
		}
	} else if subCommand == "jobdef" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\tjobdef <job_name>")
		}

		jobName := flag.Arg(1)
		err := controller.PrintJobDefinition(jobName)
		if err != nil {
			return err
		}
	} else if subCommand == "jobhist" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\tjobhist <job_name>")
		}

		jobName := flag.Arg(1)
		controller.PrintJobHistory(jobName)
	}

	return nil
}

// dependents is a subcommand to view job relations
func dependents(subcommand string, client pb.JobDependsClient) error {
	return nil
}

// printHelp prints subcommand usage
func printHelp() {
	helpStr := `Usage:
	submitjil <file>
	sendevent <job_name> <event_type>
	jobstat <job_name>
	jobdef <job_name>
	jobhist <job_name>
	jobdepend <job_name>
	jobfuture <job_name>`

	fmt.Println(helpStr)
}

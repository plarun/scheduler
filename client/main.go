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

const (
	sendEventUsage = "event <job_name> <event_type = start,abort,freeze,reset,green>"
	submitJilUsage = "submit <file_path>"
	statusUsage    = "status <job_name>"
	jobUsage       = "job <job_name>"
	historyUsage   = "history <job_name>"
	// assocUsage     = "assoc <job_name>"
	// futureUsage    = "future <job_name>"
)

func main() {
	startClient()
}

// startClient starts and creates all required client services.
func startClient() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing subcommand")
		printHelp()
		os.Exit(1)
	}

	// client conn to event-server
	conn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	cmd := flag.Arg(0)

	if cmd == "submit" {
		err = submitJil(pb.NewSubmitJilClient(conn))
	} else if cmd == "event" {
		err = sendEvent(pb.NewSendEventClient(conn))
	} else if cmd == "status" || cmd == "job" || cmd == "history" {
		err = status(cmd, pb.NewJobStatusClient(conn))
	} else if cmd == "assoc" || cmd == "future" {
		err = dependents(cmd, pb.NewJobDependsClient(conn))
	} else {
		_, err := fmt.Fprintln(os.Stderr, "invalid subcommand")
		if err != nil {
			return
		}
		printHelp()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// sendevent is a subcommand to send a job action event request
func sendEvent(client pb.SendEventClient) error {
	if flag.NArg() != 3 {
		return fmt.Errorf("invalid argument\nusage:\n\t%s", sendEventUsage)
	}

	jobName := flag.Arg(1)
	eventType := flag.Arg(2)

	controller := job.NewEventController(client)
	if err := controller.Event(jobName, eventType); err != nil {
		return err
	}

	return nil
}

// submitjil is a subcommand to insert, update or delete the job definitions
func submitJil(client pb.SubmitJilClient) error {
	if flag.NArg() != 2 {
		return fmt.Errorf("invalid argument\nusage:\n\t%s", submitJilUsage)
	}

	inputFilename := flag.Arg(1)

	controller := job.NewJobInfoController(client)
	if err := controller.SubmitJil(inputFilename); err != nil {
		return err
	}

	return nil
}

// status is a subcommand to view the latest runs, jobs definition and run history
func status(subCommand string, client pb.JobStatusClient) error {
	controller := job.NewJobStatusController(client)

	if subCommand == "status" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\t%s", statusUsage)
		}

		jobName := flag.Arg(1)

		if err := controller.PrintJobStatus(jobName); err != nil {
			return err
		}
	} else if subCommand == "job" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\t%s", jobUsage)
		}

		jobName := flag.Arg(1)

		if err := controller.PrintJobDefinition(jobName); err != nil {
			return err
		}
	} else if subCommand == "history" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\t%s", historyUsage)
		}

		jobName := flag.Arg(1)

		if err := controller.PrintJobHistory(jobName); err != nil {
			return err
		}
	}

	return nil
}

// dependents is a subcommand to view job relations
func dependents(subcommand string, client pb.JobDependsClient) error {
	return nil
}

// printHelp prints subcommand usage
func printHelp() {
	fmt.Printf("Usage:\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n",
		submitJilUsage,
		sendEventUsage,
		statusUsage,
		jobUsage,
		historyUsage)
}

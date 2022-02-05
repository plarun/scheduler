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
		fmt.Fprintln(os.Stderr, "invalid subcommand")
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
		return fmt.Errorf("invalid argument\nusage:\n\tevent <job_name> <event_type = start,abort,freeze,reset,green>")
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
		return fmt.Errorf("invalid argument\nusage:\n\tsubmit <file_path>")
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

	if subCommand == "status" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\tstatus <job_name>")
		}

		jobName := flag.Arg(1)

		if err := controller.PrintJobStatus(jobName); err != nil {
			return err
		}
	} else if subCommand == "job" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\tjob <job_name>")
		}

		jobName := flag.Arg(1)

		if err := controller.PrintJobDefinition(jobName); err != nil {
			return err
		}
	} else if subCommand == "history" {
		if flag.NArg() != 2 {
			return fmt.Errorf("invalid argument\nusage:\n\thistory <job_name>")
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
	helpStr := `Usage:
	submit <file>
	event <job_name> <event_type = start,abort,freeze,reset,green>
	status <job_name>
	job <job_name>
	history <job_name>
	assoc <job_name>
	future <job_name>`

	fmt.Println(helpStr)
}

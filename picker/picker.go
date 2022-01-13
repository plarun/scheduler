package main

import (
	"log"
	"time"

	pb "github.com/plarun/scheduler/picker/data"
	"golang.org/x/net/context"
)

// JobPicker wraps the NextJobsClient and queues the next run jobs
type JobPicker struct {
	client pb.NextJobsClient
}

func NewJobPicker(client pb.NextJobsClient) *JobPicker {
	return &JobPicker{client: client}
}

// Pick get and pushes the next run jobs to waiting queue
func (picker JobPicker) Pick() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	req := &pb.NextJobsReq{}
	res, err := picker.client.Next(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	for _, job := range res.JobList {
		waitingQueue.Push(job)
	}
}

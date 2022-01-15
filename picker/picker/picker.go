package picker

import (
	"time"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/queue"
	"golang.org/x/net/context"
)

// JobPicker wraps the NextJobsClient and queues the next run jobs
type JobPicker struct {
	Client pb.NextJobsClient
	Queue  *queue.ConcurrentWaitingQueue
}

func NewJobPicker(client pb.NextJobsClient) *JobPicker {
	return &JobPicker{
		Client: client,
		Queue:  queue.NewWaitingQueue(),
	}
}

// Pick get and pushes the next run jobs to waiting queue
func (picker JobPicker) NextJobs() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	nextJobReq := &pb.NextJobsReq{}
	nextJobRes, err := picker.Client.Next(ctx, nextJobReq)
	if err != nil {
		return err
	}

	for _, job := range nextJobRes.JobList {
		picker.Queue.Push(job)
	}

	return nil
}

package pickpass

import (
	"time"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/wait"
	"golang.org/x/net/context"
)

// JobPicker wraps the NextJobsClient and queues the next run jobs
type JobPicker struct {
	PickClient pb.PickJobsClient
	PassClient pb.PassJobsClient
	Queue      *wait.ConcurrentWaitingQueue
	Holder     *wait.ConcurrentHolder
}

func NewJobPicker(pickClient pb.PickJobsClient, passClient pb.PassJobsClient) *JobPicker {
	return &JobPicker{
		PickClient: pickClient,
		PassClient: passClient,
		Queue:      wait.NewWaitingQueue(),
		Holder:     wait.NewConcurrentHolder(),
	}
}

// NextJobs get and pushes the next run jobs to waiting queue
func (picker JobPicker) PickJobs() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pickJobReq := &pb.PickJobsReq{}
	pickJobRes, err := picker.PickClient.Pick(ctx, pickJobReq)
	if err != nil {
		return err
	}

	for _, job := range pickJobRes.JobList {
		if job.ConditionSatisfied {
			picker.Queue.Push(job)
		} else {
			picker.Holder.Hold(job)
		}
	}

	return nil
}

// PassJobs passes the jobs in queue to controller
func (picker JobPicker) PassJobs() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if picker.Queue.Size() == 0 {
		return nil
	}

	var jobList []*pb.Job

	for picker.Queue.Size() != 0 {
		job, err := picker.Queue.Pop()
		if err != nil {
			return err
		}
		jobList = append(jobList, job.Job())
	}

	passJobReq := &pb.PassJobsReq{
		JobList: jobList,
	}
	_, err := picker.PassClient.Pass(ctx, passJobReq)
	if err != nil {
		return err
	}

	return nil
}

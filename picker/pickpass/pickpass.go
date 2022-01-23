package pickpass

import (
	"time"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/wait"
	"golang.org/x/net/context"
)

var PickPass *JobPicker = nil

// JobPicker wraps the NextJobsClient and queues the next run jobs
type JobPicker struct {
	PickClient pb.PickJobsClient
	PassClient pb.PassJobsClient
	Holder     *wait.ConcurrentHolder
}

func GetPickPass(pickClient pb.PickJobsClient, passClient pb.PassJobsClient) *JobPicker {
	if PickPass == nil {
		PickPass = &JobPicker{
			PickClient: pickClient,
			PassClient: passClient,
			Holder:     wait.NewConcurrentHolder(),
		}
	}

	return PickPass
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
			PassJobs(job)
		} else {
			picker.Holder.Hold(job)
		}
	}

	return nil
}

// PassJobs passes the jobs in queue to controller
func PassJobs(job *pb.Job) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Already initialized
	picker := GetPickPass(nil, nil)

	passJobReq := &pb.PassJobsReq{
		ReadyJob: job,
	}

	_, err := picker.PassClient.Pass(ctx, passJobReq)
	if err != nil {
		return err
	}

	return nil
}

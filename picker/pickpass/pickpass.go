package pickpass

import (
	"log"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/wait"
	"golang.org/x/net/context"
)

// Singleton JobPicker instance
var pickPass *JobPicker = nil

// JobPicker wraps the NextJobsClient and queues the next run jobs
type JobPicker struct {
	PickClient pb.PickJobsClient
	PassClient pb.PassJobsClient
	Holder     *wait.ConcurrentHolder
}

// InitPickPass initiates the JobPicker
func InitPickPass(pickClient pb.PickJobsClient, passClient pb.PassJobsClient) *JobPicker {
	if pickPass == nil {
		pickPass = &JobPicker{
			PickClient: pickClient,
			PassClient: passClient,
			Holder:     wait.NewConcurrentHolder(),
		}
	}

	return pickPass
}

// GetPickPass returns the already initiated singleton JobPicker
func GetPickPass() *JobPicker {
	return pickPass
}

// NextJobs get and pushes the next run jobs to waiting queue
func (picker JobPicker) PickJobs() error {
	pickJobReq := &pb.PickJobsReq{}
	pickJobRes, err := picker.PickClient.Pick(context.Background(), pickJobReq)
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

	picker.Holder.Print()
	return nil
}

// PassJobs passes the jobs in queue to controller
func PassJobs(job *pb.ReadyJob) error {
	// Already initialized
	picker := GetPickPass()
	passJobReq := &pb.PassJobsReq{
		Job: job,
	}

	log.Printf("Job: %s is passed to controller\n", job.GetJobName())

	_, err := picker.PassClient.Pass(context.Background(), passJobReq)
	if err != nil {
		return err
	}

	return nil
}

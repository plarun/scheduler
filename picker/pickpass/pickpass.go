package pickpass

import (
	"log"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/wait"
	"golang.org/x/net/context"
)

// Singleton JobPickPasser instance
var pickPass *JobPickPasser = nil

// JobPickPasser wraps the NextJobsClient and queues the next run jobs
type JobPickPasser struct {
	PickClient pb.PickJobsClient
	PassClient pb.PassJobsClient
	Holder     *wait.ConcurrentHolder
}

// InitPickPasser initiates the JobPickPasser
func InitPickPasser(pickClient pb.PickJobsClient, passClient pb.PassJobsClient) *JobPickPasser {
	if pickPass == nil {
		pickPass = &JobPickPasser{
			PickClient: pickClient,
			PassClient: passClient,
			Holder:     wait.NewConcurrentHolder(),
		}
	}

	return pickPass
}

func GetPickPasser() *JobPickPasser {
	return pickPass
}

// PickJobs get and pushes the next run jobs to waiting queue
func (picker *JobPickPasser) PickJobs() error {
	pickJobReq := &pb.PickJobsReq{}
	pickJobRes, err := picker.PickClient.Pick(context.Background(), pickJobReq)
	if err != nil {
		return err
	}

	for _, job := range pickJobRes.JobList {
		if job.ConditionSatisfied {
			err := picker.PassJobs(job)
			if err != nil {
				return err
			}
		} else {
			picker.Holder.Hold(job)
		}
	}

	picker.Holder.Print()
	return nil
}

// PassJobs passes the jobs in queue to controller
func (picker *JobPickPasser) PassJobs(job *pb.ReadyJob) error {
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

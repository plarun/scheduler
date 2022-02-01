package job

import (
	"log"
	"time"

	pb "github.com/plarun/scheduler/client/data"
	"golang.org/x/net/context"
)

type JobStatusController struct {
	client pb.JobStatusClient
}

func NewJobStatusController(client pb.JobStatusClient) *JobStatusController {
	return &JobStatusController{
		client: client,
	}
}

func (controller JobStatusController) PrintJobStatus(jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	jobStatusReq := &pb.GetJobRunStatusReq{
		JobName: jobName,
	}

	jobStatusRes, err := controller.client.GetJobRunStatus(ctx, jobStatusReq)
	if err != nil {
		return err
	}

	runStatus(
		jobStatusRes.GetJobName(),
		jobStatusRes.GetStartTime(),
		jobStatusRes.GetEndTime(),
		jobStatusRes.GetStatusType().String())

	return nil
}

func (controller JobStatusController) PrintJobDefinition(jobName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	jobDefinitionReq := &pb.GetJilReq{
		JobName: jobName,
	}

	log.Println(jobDefinitionReq.JobName)

	jobDefinitionRes, err := controller.client.GetJobDefinition(ctx, jobDefinitionReq)
	if err != nil {
		return err
	}

	jobDefinition(jobDefinitionRes)

	return nil
}

func (controller JobStatusController) PrintJobHistory(jobName string) error {
	return nil
}

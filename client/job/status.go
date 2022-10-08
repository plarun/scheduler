package job

import (
	"fmt"

	pb "github.com/plarun/scheduler/client/proto"
	"golang.org/x/net/context"
)

type StatusController struct {
	client pb.JobStatusClient
}

// NewJobStatusController returns new instance of JobStatusController
func NewJobStatusController(client pb.JobStatusClient) *StatusController {
	return &StatusController{
		client: client,
	}
}

// PrintJobStatus controls on printing the last run status of job
func (controller StatusController) PrintJobStatus(jobName string) error {
	jobStatusReq := &pb.GetJobRunStatusReq{
		JobName: jobName,
	}

	jobStatusRes, err := controller.client.GetJobRunStatus(context.Background(), jobStatusReq)
	if err != nil {
		return fmt.Errorf("job not found")
	}

	printRunStatus(
		jobStatusRes.GetJobName(),
		jobStatusRes.GetStartTime(),
		jobStatusRes.GetEndTime(),
		jobStatusRes.GetStatusType().String())
	return nil
}

// PrintJobDefinition controls on printing the job definition
func (controller StatusController) PrintJobDefinition(jobName string) error {
	jobDefinitionReq := &pb.GetJilReq{
		JobName: jobName,
	}

	jobDefinitionRes, err := controller.client.GetJobDefinition(context.Background(), jobDefinitionReq)
	if err != nil {
		return fmt.Errorf("job not found")
	}

	printJobDefinition(jobDefinitionRes)
	return nil
}

// PrintJobHistory controls on printing previous run history of job
func (controller StatusController) PrintJobHistory(jobName string) error {
	jobRunHistoryReq := &pb.GetJobRunHistoryReq{
		JobName: jobName,
	}

	jobRunHistoryRes, err := controller.client.GetJobRunHistory(context.Background(), jobRunHistoryReq)
	if err != nil {
		return fmt.Errorf("job not found")
	}

	printRunHistory(jobName, jobRunHistoryRes.StartTime, jobRunHistoryRes.EndTime, jobRunHistoryRes.StatusType)
	return nil
}

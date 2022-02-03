package job

import (
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
	ctx := context.Background()

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
	ctx := context.Background()

	jobDefinitionReq := &pb.GetJilReq{
		JobName: jobName,
	}

	jobDefinitionRes, err := controller.client.GetJobDefinition(ctx, jobDefinitionReq)
	if err != nil {
		return err
	}

	jobDefinition(jobDefinitionRes)

	return nil
}

func (controller JobStatusController) PrintJobHistory(jobName string) error {
	ctx := context.Background()

	jobRunHistoryReq := &pb.GetJobRunHistoryReq{
		JobName: jobName,
	}

	jobRunHistoryRes, err := controller.client.GetJobRunHistory(ctx, jobRunHistoryReq)
	if err != nil {
		return err
	}

	runHistory(jobName, jobRunHistoryRes.StartTime, jobRunHistoryRes.EndTime, jobRunHistoryRes.StatusType)

	return nil
}

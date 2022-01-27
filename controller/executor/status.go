package executor

import (
	"context"
	"time"

	pb "github.com/plarun/scheduler/controller/data"
)

var statusClient *UpdateStatusClient = nil

type UpdateStatusClient struct {
	client pb.UpdateStatusClient
}

func InitUpdateStatusClient(client pb.UpdateStatusClient) {
	statusClient = &UpdateStatusClient{
		client: client,
	}
}

func updateStatus(jobName string, status pb.NewStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	updateJobStatusReq := &pb.UpdateStatusReq{
		JobName: jobName,
		Status:  status,
	}
	_, err := statusClient.client.Update(ctx, updateJobStatusReq)
	if err != nil {
		return err
	}
	return nil
}

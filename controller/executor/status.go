package executor

import (
	"context"
	"fmt"
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
	if _, err := statusClient.client.Update(ctx, updateJobStatusReq); err != nil {
		return fmt.Errorf("updateStatus: %v", err)
	}

	return nil
}

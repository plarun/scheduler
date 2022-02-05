package executor

import (
	"context"
	"fmt"

	pb "github.com/plarun/scheduler/controller/data"
)

var statusClient *UpdateStatusClient = nil

type UpdateStatusClient struct {
	client pb.UpdateStatusClient
}

// InitUpdateStatusClient initiates a new UpdateStatusClient
func InitUpdateStatusClient(client pb.UpdateStatusClient) {
	statusClient = &UpdateStatusClient{
		client: client,
	}
}

// updateStatus requests the monitor for job status update
func updateStatus(jobName string, status pb.NewStatus) error {
	updateJobStatusReq := &pb.UpdateStatusReq{
		JobName: jobName,
		Status:  status,
	}

	if _, err := statusClient.client.Update(context.Background(), updateJobStatusReq); err != nil {
		return fmt.Errorf("updateStatus: %v", err)
	}

	return nil
}

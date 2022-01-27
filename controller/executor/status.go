package executor

import (
	pb "github.com/plarun/scheduler/controller/data"
)

var statusClient *UpdateStatusClient = nil

type UpdateStatusClient struct {
	client *pb.UpdateStatusClient
}

func InitUpdateStatusClient(client *pb.UpdateStatusClient) {
	statusClient = &UpdateStatusClient{
		client: client,
	}
}

func getUpdateStatusClient() *UpdateStatusClient {
	return statusClient
}

func updateStatus(jobName string) error {

	return nil
}

package model

import (
	pb "github.com/plarun/scheduler/client/data"
)

type ClientServices struct {
	SubmitJil pb.SubmitJilClient
	SendEvent pb.SendEventClient
	JobStatus pb.JobStatusClient
	Dependent pb.JobDependsClient
}

func NewClientServices() *ClientServices {
	return &ClientServices{}
}

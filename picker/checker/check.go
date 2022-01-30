package checker

import (
	"context"
	"log"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/pickpass"
	"github.com/plarun/scheduler/picker/wait"
)

type HoldChecker struct {
	pb.UnimplementedConditionServer
	Holder *wait.ConcurrentHolder
}

func NewHoldChecker() *HoldChecker {
	return &HoldChecker{
		Holder: wait.NewConcurrentHolder(),
	}
}

func (checker HoldChecker) ConditionStatus(ctx context.Context, req *pb.JobConditionReq) (*pb.JobConditionRes, error) {
	log.Printf("Satisfied Successors: %v\n", req.GetSatisfiedSuccessors())
	for _, dependentJob := range req.GetSatisfiedSuccessors() {
		if checker.Holder.Contains(dependentJob) {
			job := checker.Holder.Free(dependentJob)
			pickpass.PassJobs(job)

		}
	}

	return &pb.JobConditionRes{}, nil
}

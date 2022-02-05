package checker

import (
	"context"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/pickpass"
	"github.com/plarun/scheduler/picker/wait"
)

type HoldChecker struct {
	pb.UnimplementedConditionServer
	Holder *wait.ConcurrentHolder
}

// NewHolderChecker returns new instance of HolderChecker
func NewHoldChecker() *HoldChecker {
	return &HoldChecker{
		Holder: wait.NewConcurrentHolder(),
	}
}

// ConditionStatus checks on the successors of the successfully completed job
func (checker HoldChecker) ConditionStatus(ctx context.Context, req *pb.JobConditionReq) (*pb.JobConditionRes, error) {
	for _, dependentJob := range req.GetSatisfiedSuccessors() {
		if checker.Holder.Contains(dependentJob) {
			job := checker.Holder.Free(dependentJob)
			pickpass.PassJobs(job)
		}
	}

	return &pb.JobConditionRes{}, nil
}

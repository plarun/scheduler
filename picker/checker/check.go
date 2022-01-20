package checker

import (
	"context"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/wait"
)

type HoldChecker struct {
	pb.UnimplementedConditionServer
	Holder *wait.ConcurrentHolder
	Queue  *wait.ConcurrentWaitingQueue
}

func NewHoldChecker() *HoldChecker {
	return &HoldChecker{
		Holder: wait.NewConcurrentHolder(),
		Queue:  wait.NewWaitingQueue(),
	}
}

func (checker HoldChecker) ConditionStatus(ctx context.Context, req *pb.JobConditionReq) (*pb.JobConditionRes, error) {
	if req.GetConditionSatisfied() && checker.Holder.Contains(req.GetJobName()) {
		job := &pb.Job{
			JobName:            req.GetJobName(),
			ConditionSatisfied: true,
		}
		checker.Queue.Push(job)
	}

	return &pb.JobConditionRes{}, nil
}

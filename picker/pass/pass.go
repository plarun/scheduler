package pass

import (
	"context"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/wait"
)

type JobPass struct {
	pb.UnimplementedPassJobsServer
	Queue *wait.ConcurrentWaitingQueue
}

func NewJobPass() *JobPass {
	return &JobPass{
		Queue: wait.NewWaitingQueue(),
	}
}

func (pass JobPass) Pass(ctx context.Context, req *pb.PassJobsReq) (*pb.PassJobsRes, error) {
	var passList []*pb.ProcessJob

	// todo: pass jobs as stream, meanwhile check job's condition

	res := &pb.PassJobsRes{
		JobList: passList,
	}
	return res, nil
}

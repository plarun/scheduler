package pass

import (
	"context"

	pb "github.com/plarun/scheduler/picker/data"
	"github.com/plarun/scheduler/picker/queue"
)

type JobPass struct {
	pb.UnimplementedPassJobsServer
	Queue *queue.ConcurrentWaitingQueue
}

func NewJobPass() *JobPass {
	return &JobPass{
		Queue: queue.NewWaitingQueue(),
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

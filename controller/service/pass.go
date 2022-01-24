package receiver

import (
	"context"

	pb "github.com/plarun/scheduler/controller/data"
	"github.com/plarun/scheduler/controller/queue"
)

type PassJobsServer struct {
	pb.UnimplementedPassJobsServer
	queue *queue.ConcurrentProcessQueue
}

func NewPassJobsServer() *PassJobsServer {
	return &PassJobsServer{
		queue: queue.GetProcessQueue(),
	}
}

func (pass PassJobsServer) Pass(ctx context.Context, req *pb.PassJobsReq) (*pb.PassJobsRes, error) {
	job := req.GetReadyJob()
	pass.queue.Push(job)

	return &pb.PassJobsRes{}, nil
}

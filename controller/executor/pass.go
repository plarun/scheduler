package executor

import (
	"context"
	"log"

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
	updateStatus(job.GetJobName())
	pass.queue.Push(job)

	log.Printf("Pass service")
	pass.queue.Print()

	return &pb.PassJobsRes{}, nil
}

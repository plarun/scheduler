package locker

import (
	"context"

	pb "github.com/plarun/scheduler/monitor/data"
)

type CheckLockServer struct {
	pb.UnimplementedCheckLockServer
	jobLocker *Locker
}

// NewCheckLockServer returns new instance of CheckLockServer
func NewCheckLockServer() *CheckLockServer {
	return &CheckLockServer{
		jobLocker: GetLocker(),
	}
}

// Check checks whether the job is locked in locker or not
func (locker CheckLockServer) Check(ctx context.Context, req *pb.CheckLockReq) (*pb.CheckLockRes, error) {
	jobName := req.GetJobName()
	locked, err := locker.jobLocker.Locked(jobName)
	if err != nil {
		return &pb.CheckLockRes{}, err
	}

	checkLockRes := &pb.CheckLockRes{
		JobName: jobName,
		Locked:  locked,
	}

	return checkLockRes, nil
}

package service

import (
	"context"
	"fmt"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
)

// ExitCodeServer represents the exit code update on job
type ExitCodeServer struct {
	pb.UnimplementedRunStatusServer
	Database *query.Database
}

// Update updates the status of job by exitcode from controller
func (excode ExitCodeServer) Update(ctx context.Context, req *pb.RunStatusReq) (*pb.RunStatusRes, error) {
	jobName := req.GetJobName()
	exitCode := req.GetStatus()

	dbTxn, err := excode.Database.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			dbTxn.Rollback()
		}
		dbTxn.Commit()
	}()

	var status pb.Status
	if exitCode == pb.ExitCode_SU {
		status = pb.Status_SUCCESS
	} else if exitCode == pb.ExitCode_FA {
		status = pb.Status_FAILED
	} else if exitCode == pb.ExitCode_AB {
		status = pb.Status_ABORTED
	} else {
		return nil, fmt.Errorf("invalid exit code type")
	}

	err = excode.Database.ChangeStatus(dbTxn, jobName, status)
	if err != nil {
		return nil, err
	}

	return &pb.RunStatusRes{}, nil
}

package service

import (
	"context"
	"fmt"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
)

// UpdateStatusServer represents the status update on job
type UpdateStatusServer struct {
	pb.UnimplementedUpdateStatusServer
	Database *query.Database
}

// Update updates the status of job by exitcode from controller
func (updStatus UpdateStatusServer) Update(ctx context.Context, req *pb.UpdateStatusReq) (*pb.UpdateStatusRes, error) {
	jobName := req.GetJobName()
	exitCode := req.GetStatus()

	dbTxn, err := updStatus.Database.DB.Begin()
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
	if exitCode == pb.NewStatus_CHANGE_SUCCESS {
		status = pb.Status_SUCCESS
	} else if exitCode == pb.NewStatus_CHANGE_FAILED {
		status = pb.Status_FAILED
	} else if exitCode == pb.NewStatus_CHANGE_ABORTED {
		status = pb.Status_ABORTED
	} else if exitCode == pb.NewStatus_CHANGE_READY {
		status = pb.Status_READY
	} else if exitCode == pb.NewStatus_CHANGE_RUNNING {
		status = pb.Status_RUNNING
	} else {
		return nil, fmt.Errorf("invalid exit code type")
	}

	err = updStatus.Database.ChangeStatus(dbTxn, jobName, status)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateStatusRes{}, nil
}

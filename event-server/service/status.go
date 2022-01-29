package service

import (
	"context"
	"fmt"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
	"google.golang.org/grpc"
)

var updateStatusService *UpdateStatusService = nil

// UpdateStatusServer represents the status update on job
type UpdateStatusService struct {
	pb.UnimplementedUpdateStatusServer
	Database      *query.Database
	MonitorClient pb.ConditionClient
}

func InitUpdateStatusService(monitorClient *grpc.ClientConn) {
	updateStatusService = &UpdateStatusService{
		Database:      query.GetDatabase(),
		MonitorClient: pb.NewConditionClient(monitorClient),
	}
}

func GetUpdateStatusService() *UpdateStatusService {
	return updateStatusService
}

// Update updates the status of job by exitcode from controller
func (updStatus UpdateStatusService) Update(ctx context.Context, req *pb.UpdateStatusReq) (*pb.UpdateStatusRes, error) {
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

	if err = updStatus.Database.ChangeStatus(dbTxn, jobName, status); err != nil {
		return nil, err
	}

	if status == pb.Status_SUCCESS {
		jobSeqId, err := updStatus.Database.GetJobId(dbTxn, jobName)
		if err != nil {
			return &pb.UpdateStatusRes{}, err
		}
		successors, err := updStatus.Database.GetSuccessors(dbTxn, jobSeqId)
		if err != nil {
			return &pb.UpdateStatusRes{}, err
		}

		conditionReq := &pb.JobConditionReq{
			JobName:    jobName,
			Successors: successors,
		}

		if _, err := updStatus.MonitorClient.ConditionStatus(ctx, conditionReq); err != nil {
			return &pb.UpdateStatusRes{}, err
		}
	}

	return &pb.UpdateStatusRes{}, nil
}

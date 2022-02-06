package service

import (
	"context"
	"fmt"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/model"
	"github.com/plarun/scheduler/event-server/query"
)

type StatusServer struct {
	Database *query.Database
	pb.UnimplementedJobStatusServer
}

// GetJobRunStatus gets the latest run status of job
func (server StatusServer) GetJobRunStatus(ctx context.Context, req *pb.GetJobRunStatusReq) (*pb.GetJobRunStatusRes, error) {
	jobName := req.GetJobName()

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer dbTxn.Commit()

	startTime, endTime, status, err := server.Database.LastRun(dbTxn, jobName)
	if err != nil {
		return nil, err
	}

	res := &pb.GetJobRunStatusRes{
		JobName:    jobName,
		StartTime:  startTime,
		EndTime:    endTime,
		StatusType: model.StatusTypeConv[status],
	}

	return res, nil
}

// GetJobDefinition gets the job definition
func (server StatusServer) GetJobDefinition(ctx context.Context, req *pb.GetJilReq) (*pb.GetJilRes, error) {
	jobName := req.GetJobName()

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return &pb.GetJilRes{}, err
	}
	defer dbTxn.Commit()

	res, err := server.Database.GetJobData(dbTxn, jobName)
	if err != nil {
		return &pb.GetJilRes{}, err
	}

	return res, nil
}

// GetJobRunHistory gets the previous runs of job
func (server StatusServer) GetJobRunHistory(ctx context.Context, req *pb.GetJobRunHistoryReq) (*pb.GetJobRunHistoryRes, error) {
	jobName := req.GetJobName()
	res := &pb.GetJobRunHistoryRes{}

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer dbTxn.Commit()

	if found := server.Database.CheckJob(jobName); !found {
		return &pb.GetJobRunHistoryRes{}, fmt.Errorf("job not found")
	}

	startTimes, endTimes, statuses, err := server.Database.GetRunHistory(dbTxn, jobName)
	if err != nil {
		return res, err
	}

	res.StartTime = startTimes
	res.EndTime = endTimes
	var convStatus []pb.Status = make([]pb.Status, 0)
	for _, status := range statuses {
		convStatus = append(convStatus, pb.Status(pb.Status_value[status]))
	}
	res.StatusType = convStatus

	return res, nil
}

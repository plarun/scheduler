package service

import (
	"context"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/model"
	"github.com/plarun/scheduler/event-server/query"
)

type StatusServer struct {
	Database *query.Database
	pb.UnimplementedJobStatusServer
}

func (server StatusServer) GetJobRunStatus(ctx context.Context, req *pb.GetJobRunStatusReq) (*pb.GetJobRunStatusRes, error) {
	jobName := req.GetJobName()

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer dbTxn.Rollback()

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

func (server StatusServer) GetJil(ctx context.Context, req *pb.GetJilReq) (*pb.GetJilRes, error) {
	jobName := req.GetJobName()

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer dbTxn.Rollback()

	res, err := server.Database.GetJobData(dbTxn, jobName)
	if err != nil {
		return nil, err
	}

	return res, nil
}

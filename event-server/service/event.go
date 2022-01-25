package service

import (
	"context"
	"database/sql"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
)

type SendEventServer struct {
	Database *query.Database
	pb.UnimplementedSendEventServer
}

func (server SendEventServer) Send(ctx context.Context, req *pb.SendEventReq) (*pb.SendEventRes, error) {
	jobName := req.GetJobName()
	eventType := req.GetEventType()
	res := &pb.SendEventRes{
		JobName:       jobName,
		EventChanged:  false,
		CurrentStatus: "",
	}

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			dbTxn.Rollback()
		}
		dbTxn.Commit()
	}()

	currStatus, err := server.Database.GetStatus(dbTxn, jobName)
	if err != nil {
		return nil, err
	}

	switch eventType {
	case pb.Event_START:
		// todo
	case pb.Event_ABORT:
		// todo
	case pb.Event_RESET:
		err = server.reset(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_FREEZE:
		err = server.freeze(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_GREEN:
		err = server.markAsSuccess(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	}

	currStatus, err = server.Database.GetStatus(dbTxn, jobName)
	if err != nil {
		return nil, err
	}
	res.CurrentStatus = currStatus.String()

	return res, nil
}

func (server SendEventServer) markAsSuccess(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_ABORTED || currStatus == pb.Status_FAILED || currStatus == pb.Status_IDLE {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_SUCCESS)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) freeze(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_ABORTED || currStatus == pb.Status_SUCCESS || currStatus == pb.Status_IDLE || currStatus == pb.Status_FAILED {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_FROZEN)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) reset(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_FROZEN {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_IDLE)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

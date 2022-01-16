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
	defer dbTxn.Rollback()

	currStatus, err := server.Database.GetStatus(dbTxn, jobName)
	if err != nil {
		return nil, err
	}

	switch eventType {
	case pb.Event_MARK_AS_SUCCESS:
		err = server.markAsSuccess(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_MARK_AS_FAILURE:
		err = server.markAsFailure(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_ON_ICE:
		err = server.onIce(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_OFF_ICE:
		err = server.offIce(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_ON_HOLD:
		err = server.onHold(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_OFF_HOLD:
		err = server.offHold(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
	case pb.Event_START:
		// todo
	case pb.Event_FORCE_START:
		// todo
	case pb.Event_KILL:
		// todo
	}

	currStatus, err = server.Database.GetStatus(dbTxn, jobName)
	if err != nil {
		return nil, err
	}
	res.CurrentStatus = currStatus.String()

	return res, nil
}

func (server SendEventServer) markAsSuccess(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_TERMINATED || currStatus == pb.Status_FAILURE || currStatus == pb.Status_INACTIVE {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_SUCCESS)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) markAsFailure(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_TERMINATED || currStatus == pb.Status_SUCCESS || currStatus == pb.Status_INACTIVE {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_FAILURE)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) onIce(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_TERMINATED || currStatus == pb.Status_SUCCESS || currStatus == pb.Status_INACTIVE || currStatus == pb.Status_FAILURE {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_ONICE)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) offIce(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_ONICE {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_INACTIVE)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) onHold(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_TERMINATED || currStatus == pb.Status_SUCCESS || currStatus == pb.Status_INACTIVE || currStatus == pb.Status_FAILURE {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_ONHOLD)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) offHold(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_ONHOLD {
		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_INACTIVE)
		if err != nil {
			return err
		}

		res.EventChanged = true
	}

	return nil
}

package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
	"google.golang.org/grpc"
)

type SendEventServer struct {
	Database *query.Database
	pb.UnimplementedSendEventServer
}

func (server SendEventServer) Event(ctx context.Context, req *pb.SendEventReq) (*pb.SendEventRes, error) {
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
		err = server.start(dbTxn, jobName, currStatus, res)
		if err != nil {
			return nil, err
		}
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

func (server SendEventServer) start(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus != pb.Status_QUEUED &&
		currStatus != pb.Status_READY &&
		currStatus != pb.Status_RUNNING {

		jobData, err := server.Database.GetJobData(dbTxn, jobName)
		if err != nil {
			return err
		}
		forceStartReq := &pb.PassJobsReq{
			ReadyJob: &pb.Job{
				JobName:            jobData.GetJobName(),
				Command:            jobData.GetCommand(),
				Machine:            jobData.GetMachine(),
				OutFile:            jobData.GetStdOut(),
				ErrFile:            jobData.GetStdErr(),
				ConditionSatisfied: true,
			},
		}

		// client connection to collector
		controllerConn, err := grpc.Dial("localhost:5557", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("connection failed: %v", err)
		}
		defer controllerConn.Close()

		controllerClient := pb.NewPassJobsClient(controllerConn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if _, err := controllerClient.Pass(ctx, forceStartReq); err != nil {
			return err
		}
		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) markAsSuccess(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_ABORTED ||
		currStatus == pb.Status_FAILED ||
		currStatus == pb.Status_IDLE {

		err := server.Database.ChangeStatus(dbTxn, jobName, pb.Status_SUCCESS)
		if err != nil {
			return err
		}
		res.EventChanged = true
	}

	return nil
}

func (server SendEventServer) freeze(dbTxn *sql.Tx, jobName string, currStatus pb.Status, res *pb.SendEventRes) error {
	if currStatus == pb.Status_ABORTED ||
		currStatus == pb.Status_SUCCESS ||
		currStatus == pb.Status_IDLE ||
		currStatus == pb.Status_FAILED {

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

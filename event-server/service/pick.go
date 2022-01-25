package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
)

type NextJobsServer struct {
	Database *query.Database
	pb.UnimplementedPickJobsServer
}

// Next gets all the jobs which are ready for next run
func (server NextJobsServer) Pick(ctx context.Context, req *pb.PickJobsReq) (*pb.PickJobsRes, error) {
	start := time.Now()
	end := start.Add(time.Second * 5)

	var startTime string = strconv.Itoa(start.Hour()) + ":" + strconv.Itoa(start.Minute()) + ":" + strconv.Itoa(start.Second())
	var endTime string = strconv.Itoa(end.Hour()) + ":" + strconv.Itoa(end.Minute()) + ":" + strconv.Itoa(end.Second())

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

	nextJobs, err := server.Database.GetNextRunJobs(dbTxn, startTime, endTime, weekDay())
	if err != nil {
		return nil, fmt.Errorf("Pick: %v", err)
	}

	res := &pb.PickJobsRes{
		JobList: nextJobs,
	}

	return res, nil
}

// weekDay returns today's weekday as required in DB
func weekDay() string {
	week := time.Now().Weekday()
	var weekConverted string
	switch week {
	case time.Sunday:
		weekConverted = "su"
	case time.Monday:
		weekConverted = "mo"
	case time.Tuesday:
		weekConverted = "tu"
	case time.Wednesday:
		weekConverted = "we"
	case time.Thursday:
		weekConverted = "th"
	case time.Friday:
		weekConverted = "fr"
	case time.Saturday:
		weekConverted = "sa"
	}

	return weekConverted
}

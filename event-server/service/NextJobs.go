package service

import (
	"context"
	"strconv"
	"time"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
)

type NextJobsServer struct {
	Database *query.Database
	pb.UnimplementedNextJobsServer
}

// Next gets all the jobs which are ready for next run
func (server NextJobsServer) Next(ctx context.Context, req *pb.NextJobsReq) (*pb.NextJobsRes, error) {
	start := time.Now()
	end := start.Add(time.Second * 5)

	var startTime string = strconv.Itoa(start.Hour()) + ":" + strconv.Itoa(start.Minute())
	var endTime string = strconv.Itoa(end.Hour()) + ":" + strconv.Itoa(end.Minute())

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer dbTxn.Rollback()

	nextJobs, err := server.Database.GetNextRunJobs(dbTxn, startTime, endTime, weekDay())
	if err != nil {
		return nil, err
	}

	res := &pb.NextJobsRes{
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

package service

import (
	"context"
	"log"
	"time"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/query"
	"google.golang.org/grpc"
)

// Update updates the status of job by exitcode from controller
func Update(jobName string) {

	database := query.GetDatabase()
	dbTxn, err := database.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer dbTxn.Rollback()

	jobSeqId, err := database.GetJobId(dbTxn, jobName)
	if err != nil {
		panic(err)
	}
	successors, err := database.GetSuccessors(dbTxn, jobSeqId)
	if err != nil {
		panic(err)
	}

	// client connection to monitor for freeing the holded successors at picker
	conn, err := grpc.Dial("localhost:5555", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req := &pb.JobConditionReq{
		JobName:       jobName,
		DependentJobs: successors,
	}
	client := pb.NewConditionClient(conn)
	_, err = client.ConditionStatus(ctx, req)
	if err != nil {
		panic(err)
	}
}

package service

import (
	"context"
	"fmt"
	"sort"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/model"
	"github.com/plarun/scheduler/event-server/query"
)

// Server for job data implements all rpc methods of JobDataServicesServer from grpc
type JilServer struct {
	Database *query.Database
	pb.UnimplementedSubmitJilServer
}

// Submit handles the slice of validated JILs sent by client
func (server JilServer) Submit(ctx context.Context, req *pb.SubmitJilReq) (*pb.SubmitJilRes, error) {
	if req == nil {
		return nil, fmt.Errorf("nil request")
	}
	var jils []*pb.Jil = req.GetJil()

	// sort the JIL list in the order of delete, insert, update
	sort.SliceStable(jils, func(i, j int) bool {
		return jils[i].Action > jils[j].Action
	})

	// validate the JIL data
	if err := server.validateJils(jils); err != nil {
		return nil, err
	}

	// JILs to be executed together as a single transaction
	var processQueue model.QueryQueue = *model.NewQueryQueue()
	for _, jil := range jils {
		processQueue.Add(jil)
	}

	// database transaction helps to execute list of queries
	// after all queries are success then commits otherwise rollbacks
	res, err := server.TransactionJobQuery(ctx, &processQueue)
	if err != nil {
		return res, err
	}

	return res, nil
}

// validateJils validates the current JIL interms of job's availability and relations
func (server JilServer) validateJils(jils []*pb.Jil) error {
	// To keep track on the jobs in the current JIL
	var jobsToProcess map[string]pb.JilAction = make(map[string]pb.JilAction)

	for _, jil := range jils {
		jobName := jil.Data.JobName
		jobAvailable := query.DB.CheckJob(jobName)

		if jil.GetAction() == pb.JilAction_INSERT && jobAvailable {
			return fmt.Errorf("job definition for %s is already available", jobName)
		} else if jil.GetAction() != pb.JilAction_INSERT && !jobAvailable {
			return fmt.Errorf("job definition for %s is not available", jobName)
		}

		jobsToProcess[jobName] = jil.GetAction()

		// check condition jobs
		if (jil.GetAction() != pb.JilAction_DELETE) && (jil.AttributeFlag&model.CONDITIONS != 0) {
			for _, conditionJob := range jil.Data.Conditions {
				// condition job is not available in DB and not going to be inserted before this definition
				if !server.Database.CheckJob(conditionJob) && jobsToProcess[jobName] != pb.JilAction_INSERT {
					return fmt.Errorf("condition job %s is not exist", conditionJob)
				}
			}
		}
	}
	return nil
}

// DB transaction executes all the queries then commits or nothing.
func (server JilServer) TransactionJobQuery(ctx context.Context, queries *model.QueryQueue) (*pb.SubmitJilRes, error) {
	var inserted, updated, deleted int32 = 0, 0, 0
	res := &pb.SubmitJilRes{
		Created: 0,
		Updated: 0,
		Deleted: 0,
	}

	dbTxn, err := server.Database.DB.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer dbTxn.Rollback()

	for queries.HasNext() {
		query := queries.Next()
		if query.Action == pb.JilAction_INSERT {
			if err := server.Database.InsertJob(dbTxn, query); err != nil {
				return res, err
			}
			inserted++
		} else if query.Action == pb.JilAction_UPDATE {
			if err := server.Database.UpdateJob(dbTxn, query); err != nil {
				return res, err
			}
			updated++
		} else if query.Action == pb.JilAction_DELETE {
			if err := server.Database.DeleteJob(dbTxn, query.Data.JobName); err != nil {
				return res, err
			}
			deleted++
		}
	}

	if err := dbTxn.Commit(); err != nil {
		return res, err
	}

	res.Created = inserted
	res.Updated = updated
	res.Deleted = deleted
	return res, nil
}

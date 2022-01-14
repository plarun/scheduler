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
	res, err := server.Database.TransactionJobQuery(ctx, &processQueue)
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

// // validateInsertJil validates the input data for JIL action type insert
// func validateInsertJil(jobData *pb.JilData) error {

// 	return nil
// }

// // validateUpdateJil validates the input data for JIL action type update
// func validateUpdateJil(jobData *pb.JilData) error {

// 	return nil
// }

// // validateDeleteJil validates the input data for JIL action type delete
// func validateDeleteJil(jobData *pb.JilData) error {

// 	return nil
// }

package query

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/model"
)

// Check whether a job is available in job table
func CheckJob(db *sql.DB, jobName string) bool {
	var job string
	row := db.QueryRow("select job_name from job where job_name=?", jobName)
	err := row.Scan(&job)
	return err != sql.ErrNoRows
}

// DB transaction executes all the queries then commits or nothing.
func TransactionJobQuery(ctx context.Context, db *sql.DB, queries *model.QueryQueue) (*pb.SubmitJilRes, error) {

	var inserted, updated, deleted int32 = 0, 0, 0
	res := &pb.SubmitJilRes{
		Created: 0,
		Updated: 0,
		Deleted: 0,
	}

	dbTxn, err := db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer dbTxn.Rollback()

	for queries.HasNext() {
		query := queries.Next()
		if query.Action == pb.JilAction_INSERT {
			if err := InsertJob(dbTxn, query); err != nil {
				return res, err
			}
			inserted++
		} else if query.Action == pb.JilAction_UPDATE {
			if err := UpdateJob(dbTxn, query); err != nil {
				return res, err
			}
			updated++
		} else if query.Action == pb.JilAction_DELETE {
			if err := DeleteJob(dbTxn, query.Data.JobName); err != nil {
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

// InsertJob inserts a new job definition into job table
func InsertJob(dbTxn *sql.Tx, jobData *pb.Jil) error {
	log.Println(jobData.Data)
	result, err := dbTxn.Exec(
		"insert into job "+
			"(job_name,command,std_out_log,std_err_log,machine,start_times,run_days,status) "+
			"values (?,?,?,?,?,?,?,?)",
		jobData.Data.JobName,
		jobData.Data.Command,
		jobData.Data.StdOut,
		jobData.Data.StdErr,
		jobData.Data.Machine,
		jobData.Data.StartTimes,
		jobData.Data.RunDays,
		"INACTIVE",
	)
	if err != nil {
		return fmt.Errorf("insertJob: %v", err)
	}

	jobSeqId, _ := result.LastInsertId()
	dependentJobSeqIds, err := GetJobIdList(dbTxn, jobData.Data.Conditions)
	if err != nil {
		return err
	}

	if jobData.AttributeFlag&model.CONDITIONS != 0 {
		if err := InsertJobDependent(dbTxn, jobSeqId, dependentJobSeqIds); err != nil {
			return err
		}
	}

	return nil
}

// GetJobId gets job sequence ID by job name
func GetJobId(dbTxn *sql.Tx, jobName string) (int64, error) {
	var jobSeqId int64 = 0
	row := dbTxn.QueryRow("select job_seq_id from job where job_name=?", jobName)
	if err := row.Scan(&jobSeqId); err != nil {
		return jobSeqId, err
	}
	return jobSeqId, nil
}

// GetJobIdList gets list of job sequence ID for list of jobs by job name
func GetJobIdList(dbTxn *sql.Tx, jobNameList []string) ([]int64, error) {
	var jobSeqIdList []int64 = make([]int64, 0)

	for _, jobName := range jobNameList {
		if jobSeqId, err := GetJobId(dbTxn, jobName); err != nil {
			return jobSeqIdList, err
		} else {
			jobSeqIdList = append(jobSeqIdList, jobSeqId)
		}
	}

	return jobSeqIdList, nil
}

// DeleteJob deletes an existing job definition from job table
func DeleteJob(dbTxn *sql.Tx, jobName string) error {
	var jobSeqId int64
	var err error

	jobSeqId, err = GetJobId(dbTxn, jobName)
	if err != nil {
		return err
	}

	// remove all the relations of job
	err = DeleteJobRelation(dbTxn, jobSeqId)
	if err != nil {
		return err
	}

	// remove all the run history of job
	// todo

	// remove the definition of job
	_, err = dbTxn.Exec("delete from job where job_name=?", jobName)
	if err != nil {
		return fmt.Errorf("deleteJobByJobName: %v", err)
	}

	return nil
}

// UpdateJob updates one or more columns in job table by job name
func UpdateJob(dbTxn *sql.Tx, jobData *pb.Jil) error {

	columns := buildJobUpdateQuery(jobData)
	log.Println(columns)
	if len(columns) != 0 {
		_, err := dbTxn.Exec("update job set "+columns+" where job_name=?;", jobData.Data.JobName)
		if err != nil {
			return fmt.Errorf("updateJob: %v", err)
		}
	}
	if jobData.AttributeFlag&model.CONDITIONS != 0 {
		if err := UpdateJobDependents(dbTxn, jobData.Data.JobName, jobData.Data.Conditions); err != nil {
			return err
		}
	}

	return nil
}

// // Get all jobs inside the box job from job table
// func (server *JobDataServer) GetAllJobsInBox(ctx context.Context, req *pb.GetAllJobsInBoxReq) (*pb.GetAllJobsInBoxRes, error) {
// 	res := &pb.GetAllJobsInBoxRes{}
// 	var jobType string
// 	row := server.db.QueryRow("select job_type from job where job_name=?", req.JobName)
// 	err := row.Scan(&jobType)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return res, fmt.Errorf("job not found")
// 		}
// 		return res, fmt.Errorf("getByJobName%s: %v", req.JobName, err)
// 	}
// 	if jobType != "BOX" {
// 		return res, fmt.Errorf("%s is not box job", req.JobName)
// 	}

// 	rows, err := server.db.Query("select * from job where box_job_seq_id=?")
// 	if err != nil {
// 		return nil, fmt.Errorf("getAllJobsInBox: %v", err)
// 	}
// 	defer rows.Close()

// 	var childJobs []*pb.JobData
// 	for rows.Next() {
// 		childJob := &pb.JobData{}

// 		err := rows.Scan(
// 			&childJob.JobSeqId,
// 			&childJob.JobName,
// 			&childJob.Command,
// 			&childJob.StdOutLog,
// 			&childJob.StdOutErr,
// 			&childJob.Machine,
// 			&childJob.StartTime,
// 			&childJob.RunWindow,
// 			&childJob.RunDays,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("getAllJobsInBox: %v", err)
// 		}

// 		childJobs = append(childJobs, childJob)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("getAllJobsInBox: %v", err)
// 	}

// 	res.Jobs = childJobs
// 	return res, nil
// }

// Get an existing job from job table
// func (server JilServer) getJob(jobName string) (error) {
// 	res := &pb.GetJobRes{}
// 	job := res.JobData

// 	row := server.DB.QueryRow("select * from job where job_name=?", jobName)
// 	err := row.Scan(
// 		&job.JobSeqId,
// 		&job.JobName,
// 		&job.Command,
// 		&job.StdOutLog,
// 		&job.StdOutErr,
// 		&job.Machine,
// 		&job.StartTime,
// 		&job.RunWindow,
// 		&job.RunDays,
// 	)
// 	res.JobData = job
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return res, fmt.Errorf("job not found")
// 		}
// 		return res, fmt.Errorf("getByJobName%s: %v", req.JobName, err)
// 	}
// 	return res, nil
// }

package query

import (
	"database/sql"
	"fmt"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/model"
)

const (
	timeGap = "00:00:05"
)

// Check whether a job is available in job table
func (database *Database) CheckJob(jobName string) bool {
	var job string

	database.lock.Lock()
	row := database.DB.QueryRow(
		`select job_name from job 
		where job_name=?`,
		jobName)
	database.lock.Unlock()

	err := row.Scan(&job)
	return err != sql.ErrNoRows
}

// InsertJob inserts a new job definition into job table
func (database *Database) InsertJob(dbTxn *sql.Tx, jobData *pb.Jil) error {
	database.lock.Lock()
	result, err := dbTxn.Exec(
		`insert into job (job_name,command,std_out_log,std_err_log,machine,start_times,run_days,status) 
		values (?,?,?,?,?,?,?,?)`,
		jobData.Data.JobName,
		jobData.Data.Command,
		jobData.Data.StdOut,
		jobData.Data.StdErr,
		jobData.Data.Machine,
		jobData.Data.StartTimes,
		jobData.Data.RunDays,
		"INACTIVE",
	)
	database.lock.Unlock()

	if err != nil {
		return fmt.Errorf("insertJob: %v", err)
	}

	jobSeqId, _ := result.LastInsertId()
	dependentJobSeqIds, err := database.GetJobIdList(dbTxn, jobData.Data.Conditions)
	if err != nil {
		return err
	}

	if jobData.AttributeFlag&model.CONDITIONS != 0 {
		if err := database.InsertJobDependent(dbTxn, jobSeqId, dependentJobSeqIds); err != nil {
			return err
		}
	}

	return nil
}

// GetJobId gets job sequence ID by job name
func (database *Database) GetJobId(dbTxn *sql.Tx, jobName string) (int64, error) {
	var jobSeqId int64 = 0

	database.lock.Lock()
	row := dbTxn.QueryRow(
		`select job_seq_id from job 
		where job_name=?`,
		jobName)
	database.lock.Unlock()

	if err := row.Scan(&jobSeqId); err != nil {
		if err == sql.ErrNoRows {
			return jobSeqId, fmt.Errorf("job not found")
		}
		return jobSeqId, err
	}
	return jobSeqId, nil
}

// GetJobIdList gets list of job sequence ID for list of jobs by job name
func (database *Database) GetJobIdList(dbTxn *sql.Tx, jobNameList []string) ([]int64, error) {
	var jobSeqIdList []int64 = make([]int64, 0)
	for _, jobName := range jobNameList {
		if jobSeqId, err := database.GetJobId(dbTxn, jobName); err != nil {
			return jobSeqIdList, err
		} else {
			jobSeqIdList = append(jobSeqIdList, jobSeqId)
		}
	}

	return jobSeqIdList, nil
}

// DeleteJob deletes an existing job definition from job table
func (database *Database) DeleteJob(dbTxn *sql.Tx, jobName string) error {
	var jobSeqId int64
	var err error

	jobSeqId, err = database.GetJobId(dbTxn, jobName)
	if err != nil {
		return err
	}

	// remove all the relations of job
	err = database.DeleteJobRelation(dbTxn, jobSeqId)
	if err != nil {
		return err
	}

	// remove all the run history of job
	// todo

	// remove the definition of job
	database.lock.Lock()
	_, err = dbTxn.Exec(
		`delete from job 
		where job_name=?`,
		jobName)
	database.lock.Unlock()

	if err != nil {
		return fmt.Errorf("deleteJobByJobName: %v", err)
	}

	return nil
}

// UpdateJob updates one or more columns in job table by job name
func (database *Database) UpdateJob(dbTxn *sql.Tx, jobData *pb.Jil) error {
	columns := buildJobUpdateQuery(jobData)
	if len(columns) != 0 {
		database.lock.Lock()
		_, err := dbTxn.Exec(
			"update job set "+
				columns+
				" where job_name=?;",
			jobData.Data.JobName)
		database.lock.Unlock()
		if err != nil {
			return fmt.Errorf("updateJob: %v", err)
		}
	}
	if jobData.AttributeFlag&model.CONDITIONS != 0 {
		if err := database.UpdateJobDependents(dbTxn, jobData.Data.JobName, jobData.Data.Conditions); err != nil {
			return err
		}
	}

	return nil
}

// GetNextRunJobs gives list of jobs ready for next run
func (database *Database) GetNextRunJobs(dbTxn *sql.Tx, startTime string, endTime string, runDay string) ([]*pb.ReadyJob, error) {
	database.lock.Lock()
	rows, err := dbTxn.Query(
		`select job_name from job 
		where start_times between ? and ? 
		and find_in_set(?, run_days) 
		and status in ('INACTIVE', 'SUCCESS', 'FAILED', 'TERMINATED') 
		and current_time-time(last_run) > ?`,
		startTime,
		endTime,
		runDay,
		timeGap)
	database.lock.Unlock()

	if err != nil {
		return nil, err
	}

	nextJobs := make([]*pb.ReadyJob, 0)
	for rows.Next() {
		var jobName string
		if err := rows.Scan(&jobName); err != nil {
			return nil, err
		}
		job := &pb.ReadyJob{JobName: jobName}
		nextJobs = append(nextJobs, job)
	}

	return nextJobs, nil
}

// GetJobData gets job definition
func (database *Database) GetJobData(dbTxn *sql.Tx, jobName string) (*pb.GetJilRes, error) {
	res := &pb.GetJilRes{}
	var jobSeqId int64

	database.lock.Lock()
	jobRow := dbTxn.QueryRow(
		`select job_seq_id, job_name, command, std_out_log, std_err_log, machine, run_days, start_times from job
		where job_name=?`,
		jobName)
	database.lock.Unlock()

	err := jobRow.Scan(
		&jobSeqId,
		&res.JobName,
		&res.Command,
		&res.StdOut,
		&res.StdErr,
		&res.Machine,
		&res.RunDays,
		&res.StartTimes)
	if err != nil {
		return nil, err
	}

	res.Conditions, err = database.GetPreceders(dbTxn, jobSeqId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetStatus gets the current status of job
func (database *Database) GetStatus(dbTxn *sql.Tx, jobName string) (pb.Status, error) {
	var statusName string
	var status pb.Status

	database.lock.Lock()
	row := dbTxn.QueryRow(
		`select status
		from job
		where job_name=?`,
		jobName)
	database.lock.Unlock()

	err := row.Scan(&statusName)
	if err != nil {
		if err == sql.ErrNoRows {
			return status, fmt.Errorf("job not found")
		}
		return status, err
	}

	return pb.Status(pb.Status_value[statusName]), nil
}

// ChangeStatus updates the status of job
func (database *Database) ChangeStatus(dbTxn *sql.Tx, jobName string, status pb.Status) error {
	database.lock.Lock()
	_, err := dbTxn.Exec(
		`update job 
		set status=? 
		where job_name=?`,
		pb.Status_name[int32(status.Number())],
		jobName)
	database.lock.Unlock()

	if err != nil {
		return err
	}

	return nil
}

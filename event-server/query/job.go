package query

import (
	"database/sql"
	"fmt"
	"log"

	pb "github.com/plarun/scheduler/event-server/data"
	"github.com/plarun/scheduler/event-server/model"
)

// CheckJob checks whether a job is available in job table
func (database *Database) CheckJob(jobName string) bool {
	var job string
	database.lock.Lock()

	row := database.DB.QueryRow(
		`select job_name from job 
		where job_name=?`,
		jobName)

	database.lock.Unlock()

	err := row.Scan(&job)

	if database.verbose {
		log.Printf("Job: %s\n", job)
	}
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
		"IDLE",
	)

	database.lock.Unlock()

	if err != nil {
		return fmt.Errorf("insertJob: %v", err)
	}

	jobSeqId, _ := result.LastInsertId()
	dependentJobSeqIds, err := database.GetJobIdList(dbTxn, jobData.Data.Conditions)
	if err != nil {
		return fmt.Errorf("insertJob: %v", err)
	}

	if jobData.AttributeFlag&model.CONDITIONS != 0 {
		if err := database.InsertJobDependent(dbTxn, jobSeqId, dependentJobSeqIds); err != nil {
			return fmt.Errorf("insertJob: %v", err)
		}
	}

	if database.verbose {
		log.Printf("Inserted Job: %s\n", jobData.GetData().GetJobName())
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
		return jobSeqId, fmt.Errorf("getJobId: %v", err)
	}

	if database.verbose {
		log.Printf("JobSeqId of %s is %d\n", jobName, jobSeqId)
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

	if database.verbose {
		log.Printf("JobSeqIds of %v are %v\n", jobNameList, jobSeqIdList)
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
	if err := database.DeleteJobRelation(dbTxn, jobSeqId); err != nil {
		return err
	}

	// remove all the run history of job
	if err := database.ClearJobRunHistory(dbTxn, jobSeqId); err != nil {
		return err
	}

	// remove the definition of job
	database.lock.Lock()

	_, err = dbTxn.Exec(
		`delete from job 
		where job_name=?`,
		jobName)

	database.lock.Unlock()

	if err != nil {
		return fmt.Errorf("deleteJob: %v", err)
	}

	if database.verbose {
		log.Printf("Job %s is deleted\n", jobName)
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
			return fmt.Errorf("updateJob: %v", err)
		}
	}

	return nil
}

// GetNextRunJobs gives list of jobs ready for next run
func (database *Database) GetNextRunJobs(dbTxn *sql.Tx, startTime string, endTime string, runDay string) ([]*pb.ReadyJob, error) {
	database.lock.Lock()

	rows, err := dbTxn.Query(
		`select job_seq_id, job_name, command, machine, std_out_log, std_err_log
		from job
		where start_times between ? and ? 
		and find_in_set(?, run_days) 
		and status in ('IDLE', 'SUCCESS', 'FAILED', 'ABORTED')`,
		startTime,
		endTime,
		runDay)

	database.lock.Unlock()
	defer rows.Close()

	nextJobs := make([]*pb.ReadyJob, 0)

	if err != nil {
		if err == sql.ErrNoRows {
			return nextJobs, nil
		}
		return nil, fmt.Errorf("GetNextRunJobs: %v", err)
	}

	for rows.Next() {
		var jobName, command, machine, stdOut, stdErr string
		var jobSeqId int64

		if err := rows.Scan(&jobSeqId, &jobName, &command, &machine, &stdOut, &stdErr); err != nil {
			return nil, fmt.Errorf("GetNextRunJobs scanning: %v", err)
		}

		conditionSatisfied, err := database.CheckConditions(dbTxn, jobSeqId)
		if err != nil {
			return nil, fmt.Errorf("GetNextRunJobs CheckConditions: %v", err)
		}

		job := &pb.ReadyJob{
			JobName:            jobName,
			Command:            command,
			Machine:            machine,
			OutFile:            stdOut,
			ErrFile:            stdErr,
			ConditionSatisfied: conditionSatisfied,
		}
		nextJobs = append(nextJobs, job)
	}

	for _, job := range nextJobs {
		if err := database.ChangeStatus(dbTxn, job.JobName, pb.Status_QUEUED); err != nil {
			return nil, fmt.Errorf("GetNextRunJobs ChangeStatus: %v", err)
		}
	}

	if database.verbose {
		log.Printf("Time range: %s to %s on %s\n", startTime, endTime, runDay)
		log.Printf("NextJobs Count: %v\n", len(nextJobs))
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
		return &pb.GetJilRes{}, fmt.Errorf("getJobData: %v", err)
	}

	res.Conditions, err = database.GetPreceders(dbTxn, jobSeqId)
	if err != nil {
		return &pb.GetJilRes{}, fmt.Errorf("getJobData: %v", err)
	}

	if database.verbose {
		log.Printf("Job: %v\n", res.GetJobName())
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
		return status, fmt.Errorf("getStatus: %v", err)
	}

	if database.verbose {
		log.Printf("Job: %s, Status: %s\n", jobName, statusName)
	}

	return pb.Status(pb.Status_value[statusName]), nil
}

// ChangeStatus updates the status of job
func (database *Database) ChangeStatus(dbTxn *sql.Tx, jobName string, status pb.Status) error {
	statusName := pb.Status_name[int32(status.Number())]

	columns := buildJobStatusUpdateQuery(jobName, status)
	database.lock.Lock()

	_, err := dbTxn.Exec(
		"update job set "+
			columns+
			" where job_name=?",
		jobName)

	database.lock.Unlock()

	if err != nil {
		return fmt.Errorf("ChangeStatus: %v", err)
	}

	if status == pb.Status_SUCCESS || status == pb.Status_FAILED || status == pb.Status_ABORTED || status == pb.Status_FROZEN {
		if err := database.saveLastRun(dbTxn, jobName); err != nil {
			return err
		}
	}

	if database.verbose {
		log.Printf("Job: %s, UpdatedStatus: %s\n", jobName, statusName)
	}

	return nil
}

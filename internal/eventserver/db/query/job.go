package query

import (
	"database/sql"
	"fmt"
)

/*
// jobExists checks whether a job is available in database
func jobExists(jobName string) (bool, error) {
	db := db.GetDatabase()
	var isExists int

	qry := "Select Exists(Select job_name From sched_job Where job_name=?)"
	row := db.QueryRow(qry, jobName)

	err := row.Scan(&isExists)
	if err != nil && err != sql.ErrNoRows {
		return false, fmt.Errorf("jobExists: %v", err)
	}
	return isExists == 1, nil
}
*/

// getJobId gets job ID by job name
func getJobId(tx *sql.Tx, jobName string) (int64, error) {
	var jobId int64 = 0

	qry := "Select job_id From sched_job Where job_name=?"

	row := tx.QueryRow(qry, jobName)

	if err := row.Scan(&jobId); err != nil {
		if err == sql.ErrNoRows {
			return jobId, fmt.Errorf("%v job not found", jobName)
		}
		return jobId, fmt.Errorf("getJobId: %v", err)
	}

	return jobId, nil
}

// getJobIdList gets list of job ID for list of jobs by job name
func getJobIdList(tx *sql.Tx, jobNameList []string) ([]int64, error) {
	var jobIdList []int64 = make([]int64, 0)

	for _, jobName := range jobNameList {
		if jobSeqId, err := getJobId(tx, jobName); err != nil {
			return jobIdList, fmt.Errorf("getJobIdList: %v", err)
		} else {
			jobIdList = append(jobIdList, jobSeqId)
		}
	}

	return jobIdList, nil
}

/*
// GetJobData gets job definition
func GetJobData(tx *sql.Tx, jobName string) (*pb.GetJilRes, error) {
	db := db.GetDatabase()
	res := &proto.GetJilRes{}
	var jobSeqId int64

	db.RLock()
	jobRow := tx.QueryRow(cfg.Env.Query.GetJob, jobName)
	db.RUnlock()

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
		return &proto.GetJilRes{}, fmt.Errorf("GetJobData: %v", err)
	}

	res.Conditions, err = getPredecessors(tx, jobSeqId)
	if err != nil {
		return &proto.GetJilRes{}, fmt.Errorf("GetJobData: %v", err)
	}

	return res, nil
}
*/

func getRunFlag(tx *sql.Tx, jobName string) (int64, string, error) {
	var jobId int64
	var runFlag string

	qry := "Select job_id, run_flag From sched_job Where job_name=?"

	row := tx.QueryRow(qry, jobName)

	if err := row.Scan(&jobId, &runFlag); err != nil {
		if err == sql.ErrNoRows {
			return jobId, runFlag, fmt.Errorf("job not found for job_id %v", jobId)
		}
		return jobId, runFlag, fmt.Errorf("getRunFlags: %v", err)
	}
	return jobId, runFlag, nil
}

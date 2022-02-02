package query

import (
	"database/sql"
	"fmt"
)

func (database *Database) LastRun(dbTxn *sql.Tx, jobName string) (string, string, string, error) {
	var startTime, endTime, status string
	database.lock.Lock()

	row := dbTxn.QueryRow(
		`select last_start_time, last_end_time, status 
		from job
		where job_name=?`,
		jobName)

	database.lock.Unlock()

	err := row.Scan(&startTime, &endTime, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", fmt.Errorf("job not found")
		}
		return "", "", "", err
	}

	if startTime == defaultTime {
		startTime = "----:--:-- --:--:--"
	}
	if endTime == defaultTime {
		endTime = "----:--:-- --:--:--"
	}

	return startTime, endTime, status, err
}

// saveLastRun stores the last run detail
func (database *Database) saveLastRun(dbTxn *sql.Tx, jobName string) error {
	var count int8
	var lastStartTime, lastEndTime, status string

	jobSeqId, err := database.GetJobId(dbTxn, jobName)
	if err != nil {
		return fmt.Errorf("saveLastRun: %v", err)
	}

	database.lock.Lock()
	defer database.lock.Unlock()

	row := dbTxn.QueryRow(
		`select count(*) 
		from job_run_history
		where job_id=?`,
		jobSeqId)

	row2 := dbTxn.QueryRow(
		`select last_start_time, last_end_time, status
		from job
		where job_seq_id=?`,
		jobSeqId)

	if err := row.Scan(&count); err != nil {
		return fmt.Errorf("saveLastRun: %v", err)
	}

	if err := row2.Scan(&lastStartTime, &lastEndTime, &status); err != nil {
		return fmt.Errorf("saveLastRun: %v", err)
	}

	if count == 10 {
		_, err := dbTxn.Exec(
			`delete from job_run_history
			where run_seq_id = (
				select max(run_seq_id)
				from job_run_history
				where job_id = ?
			)`,
			jobSeqId)
		if err != nil {
			return fmt.Errorf("saveLastRun: %v", err)
		}
	}

	_, err = dbTxn.Exec(
		`insert into job_run_history
		(job_id, start_time, end_time, status)
		values (?,?,?,?)`,
		jobSeqId, lastStartTime, lastEndTime, status)

	if err != nil {
		return fmt.Errorf("saveLastRun: %v", err)
	}

	return nil
}

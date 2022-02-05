package query

import (
	"database/sql"
	"fmt"
	"log"
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

	row := dbTxn.QueryRow(
		`select count(*) 
		from job_run_history
		where job_id=?`,
		jobSeqId)

	if err := row.Scan(&count); err != nil {
		return fmt.Errorf("saveLastRun: %v", err)
	}

	row2 := dbTxn.QueryRow(
		`select last_start_time, last_end_time, status
		from job
		where job_seq_id=?`,
		jobSeqId)

	if err := row2.Scan(&lastStartTime, &lastEndTime, &status); err != nil {
		return fmt.Errorf("saveLastRun: %v", err)
	}

	if count == 10 {
		var runSeqId int
		row := dbTxn.QueryRow(
			`select min(run_seq_id)
			from job_run_history
			where job_id=?`,
			jobSeqId)

		if err := row.Scan(&runSeqId); err != nil {
			return fmt.Errorf("saveLastRun: %v", err)
		}

		_, err := dbTxn.Exec(
			`delete from job_run_history
			where run_seq_id = ?`,
			runSeqId)

		if err != nil {
			return fmt.Errorf("saveLastRun: %v", err)
		}
	}

	_, err = dbTxn.Exec(
		`insert into job_run_history
		(job_id, start_time, end_time, status)
		values (?,?,?,?)`,
		jobSeqId, lastStartTime, lastEndTime, status)

	database.lock.Unlock()

	if err != nil {
		return fmt.Errorf("saveLastRun: %v", err)
	}

	log.Println("exit saveLastRun")
	return nil
}

// GetRunHistory returns all the stored preivous runs statuses
func (database *Database) GetRunHistory(dbTxn *sql.Tx, jobName string) ([]string, []string, []string, error) {
	var startTimes, endTimes, statuses []string
	database.lock.Lock()

	rows, err := dbTxn.Query(
		`select start_time, end_time, status
		from job_run_history
		where job_id = (
			select job_seq_id
			from job
			where job_name = ?
		)
		order by run_seq_id desc`,
		jobName)

	database.lock.Unlock()
	if err != nil {
		return nil, nil, nil, err
	}

	var startTime, endTime, status string
	for rows.Next() {
		if err := rows.Scan(&startTime, &endTime, &status); err != nil {
			return nil, nil, nil, err
		}

		startTimes = append(startTimes, startTime)
		if endTime == defaultTime {
			endTime = "----:--:-- --:--:--"
		}
		endTimes = append(endTimes, endTime)
		statuses = append(statuses, status)
	}

	return startTimes, endTimes, statuses, nil
}

// ClearJobRunHistory removes all the run history of a given job_seq_id
func (database *Database) ClearJobRunHistory(dbTxn *sql.Tx, jobSeqId int64) error {
	database.lock.Lock()

	_, err := dbTxn.Exec(
		`delete from job_run_history
		where job_id=?`,
		jobSeqId)

	if err != nil {
		return fmt.Errorf("clearJobRunHistory: %v", err)
	}

	database.lock.Unlock()
	return nil
}

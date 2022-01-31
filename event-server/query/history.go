package query

import (
	"database/sql"
	"fmt"
)

func (database *Database) LastRun(dbTxn *sql.Tx, jobName string) (string, string, string, error) {
	var startTime, endTime, status string
	database.lock.Lock()

	rows := dbTxn.QueryRow(
		`select last_start_time, last_end_time, status 
		from job
		where job_name=?`,
		jobName)

	database.lock.Unlock()

	err := rows.Scan(&startTime, &endTime, &status)
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

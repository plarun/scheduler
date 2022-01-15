package query

import (
	"database/sql"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func (database *Database) LastRun(dbTxn *sql.Tx, jobName string) (*timestamppb.Timestamp, *timestamppb.Timestamp, string, error) {
	var startTime, endTime *timestamppb.Timestamp
	var status string

	database.lock.Lock()
	rows := dbTxn.QueryRow(
		`select last_start_time, last_end_time, status 
		from job
		where job_name=?`,
		jobName)
	database.lock.Unlock()

	err := rows.Scan(&startTime, &endTime, &status)
	if err != nil {
		return startTime, endTime, status, err
	}

	return startTime, endTime, status, nil
}

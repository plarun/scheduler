package query

import (
	"github.com/plarun/scheduler/internal/validator/db/mysql"
	"github.com/plarun/scheduler/internal/validator/errors"
)

// JobExists checks whether a job is available in job table
func JobExists(jobName string) (bool, error) {
	db := mysql.GetDatabase()
	var isExists int

	qry := "Select Exists(Select job_name From sched_job Where job_name=?)"

	row := db.DB.QueryRow(qry, jobName)

	if err := row.Scan(&isExists); err != nil {
		return false, &errors.InternalSqlError{
			Query: qry,
			Err:   err,
		}
	}
	return isExists == 1, nil
}

// GetTaskType gets the job type value of exising job
func GetTaskType(name string) (string, error) {
	database := mysql.GetDatabase()
	var typ string

	qry := "Select job_type From sched_job Where job_name=?"

	row := database.DB.QueryRow(qry, name)

	if err := row.Scan(&typ); err != nil {
		return "", &errors.InternalSqlError{
			Query: qry,
			Err:   err,
		}
	}
	return typ, nil
}

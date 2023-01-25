package query

import (
	"github.com/plarun/scheduler/internal/validator/db/mysql"
	"github.com/plarun/scheduler/internal/validator/errors"
)

func GetRunDetails(jobName string) (string, error) {
	database := mysql.GetDatabase()
	var runFlag string

	qry := "Select run_flag From sched_job Where job_name=?"

	database.RLock()
	row := database.DB.QueryRow(qry, jobName)
	database.RUnlock()

	if err := row.Scan(&runFlag); err != nil {
		return "", &errors.InternalSqlError{
			Query: qry,
			Err:   err,
		}
	}
	return runFlag, nil
}

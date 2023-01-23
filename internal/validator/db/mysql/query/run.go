package query

import (
	"github.com/plarun/scheduler/internal/validator/db/mysql"
	"github.com/plarun/scheduler/internal/validator/errors"
)

func GetRunDetails(jobName string) (int32, error) {
	database := mysql.GetDatabase()
	var runFlag int32

	qry := "Select run_flag From sched_job Where job_name=?"

	database.RLock()
	row := database.DB.QueryRow(qry, jobName)
	database.RUnlock()

	if err := row.Scan(&runFlag); err != nil {
		return 0, &errors.InternalSqlError{
			Query: qry,
			Err:   err,
		}
	}
	return runFlag, nil
}

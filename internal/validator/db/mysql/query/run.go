package query

import (
	"github.com/plarun/scheduler/internal/validator/db/mysql"
	"github.com/plarun/scheduler/internal/validator/errors"
)

func GetRunDetails(name string) (string, error) {
	database := mysql.GetDatabase()
	var runFlag string

	qry := "Select run_flag From sched_task Where name=?"

	row := database.DB.QueryRow(qry, name)

	if err := row.Scan(&runFlag); err != nil {
		return "", &errors.InternalSqlError{
			Query: qry,
			Err:   err,
		}
	}
	return runFlag, nil
}

package query

import (
	"github.com/plarun/scheduler/internal/validator/db/mysql"
	"github.com/plarun/scheduler/internal/validator/errors"
)

// TaskExists checks whether a task is already exist
func TaskExists(name string) (bool, error) {
	db := mysql.GetDatabase()
	var isExists int

	qry := "Select Exists(Select 1 From sched_task Where name=?)"

	row := db.DB.QueryRow(qry, name)

	if err := row.Scan(&isExists); err != nil {
		return false, &errors.InternalSqlError{
			Query: qry,
			Err:   err,
		}
	}
	return isExists == 1, nil
}

// GetTaskType gets the task type value of exising task
func GetTaskType(name string) (string, error) {
	database := mysql.GetDatabase()
	var typ string

	qry := "Select type From sched_task Where name=?"

	row := database.DB.QueryRow(qry, name)

	if err := row.Scan(&typ); err != nil {
		return "", &errors.InternalSqlError{
			Query: qry,
			Err:   err,
		}
	}
	return typ, nil
}

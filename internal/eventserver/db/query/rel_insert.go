package query

import (
	"database/sql"
	"fmt"

	er "github.com/plarun/scheduler/pkg/error"
)

// insertTaskRelation inserts a relation between given task id and
// given start condition task ids
func insertTaskRelation(tx *sql.Tx, id int64, condTasksIds []int64) error {
	qry := "Insert Into sched_task_relation (task_id, cond_task_id) Values (?, ?)"

	for _, condId := range condTasksIds {
		if _, err := tx.Exec(qry, id, condId); err != nil {
			return fmt.Errorf("insertTaskRelation: %w", er.NewDatabaseError(err.Error()))
		}
	}
	return nil
}

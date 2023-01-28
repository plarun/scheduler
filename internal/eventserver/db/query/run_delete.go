package query

import (
	"database/sql"
	"fmt"
)

// ClearTaskRunHistory removes all the run history of a given task id
func ClearTaskRunHistory(tx *sql.Tx, id int64) error {
	qry := "Delete From task_run_history Where task_id=?"

	_, err := tx.Exec(qry, id)

	if err != nil {
		return fmt.Errorf("ClearTaskRunHistory: %v", err)
	}

	return nil
}

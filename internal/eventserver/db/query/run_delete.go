package query

import (
	"database/sql"
	"fmt"
)

// ClearJobRunHistory removes all the run history of a given job_seq_id
func ClearJobRunHistory(tx *sql.Tx, jobId int64) error {
	qry := "Delete From job_run_history Where job_id=?"

	_, err := tx.Exec(qry, jobId)

	if err != nil {
		return fmt.Errorf("clearJobRunHistory: %v", err)
	}

	return nil
}

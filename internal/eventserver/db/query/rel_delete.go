package query

import (
	"database/sql"
	"fmt"
)

// deleteJobPredecessors removes relation between job and
// jobs in starting condition. In result there will be no
// jobs in starting condition.
func deleteJobRelation(tx *sql.Tx, jobId int64) error {
	qry := "Delete From sched_job_relation Where job_id=?"

	if _, err := tx.Exec(qry, jobId); err != nil {
		return fmt.Errorf("deleteJobPredecessors: %v", err)
	}

	return nil
}

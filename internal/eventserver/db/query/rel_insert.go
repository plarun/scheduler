package query

import (
	"database/sql"
	"fmt"
)

// insertJobRelation inserts a relation between given job id and
// given start condition job ids
func insertJobRelation(tx *sql.Tx, jobId int64, condJobIds []int64) error {
	qry := "Insert Into sched_job_relation (job_id, cond_job_id) Values (?, ?)"

	for _, condJobId := range condJobIds {
		if _, err := tx.Exec(qry, jobId, condJobId); err != nil {
			return fmt.Errorf("insertJobRelation: %v", err)
		}
	}
	return nil
}

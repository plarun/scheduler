package query

import (
	"database/sql"
	"fmt"
)

// updateJobRelation updates the relation of given job id
// with given new distnct condition jobs list. It will remove
// the existing relations of the given job id then will insert
// new relations with given new jobs list.
func updateJobRelation(dbTxn *sql.Tx, jobId int64, condJobs []string) error {
	condJobsId, err := getJobIdList(dbTxn, condJobs)
	if err != nil {
		return fmt.Errorf("updateJobRelation: %v", err)
	}

	if err := deleteJobRelation(dbTxn, jobId); err != nil {
		return fmt.Errorf("updateJobRelation: %v", err)
	}
	if err := insertJobRelation(dbTxn, jobId, condJobsId); err != nil {
		return fmt.Errorf("updateJobRelation: %v", err)
	}
	return nil
}

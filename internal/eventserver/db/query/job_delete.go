package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
)

// DeleteJob deletes an existing job definition from job table
func DeleteJob(tx *sql.Tx, tsk *task.TaskEntity) error {
	jobId, err := getJobId(tx, tsk.Name())
	if err != nil {
		return err
	}

	// remove job run times
	if err := deleteStartTimes(tx, jobId); err != nil {
		return fmt.Errorf("DeleteJob: %w", err)
	}
	if err := deleteStartMins(tx, jobId); err != nil {
		return fmt.Errorf("DeleteJob: %w", err)
	}

	// remove all the relations of job
	if err := deleteJobRelation(tx, jobId); err != nil {
		return fmt.Errorf("DeleteJob: %w", err)
	}

	// remove all the run history of job
	if err := ClearJobRunHistory(tx, jobId); err != nil {
		return fmt.Errorf("DeleteJob: %w", err)
	}

	// remove the definition of job
	qry := "Delete From sched_job Where job_name=?"
	_, err = tx.Exec(qry, tsk.Name)

	if err != nil {
		return fmt.Errorf("DeleteJob: %v", err)
	}

	return nil
}

func deleteStartTimes(tx *sql.Tx, jobId int64) error {
	qry := "Delete From sched_batch_run Where job_id=?"
	_, err := tx.Exec(qry, jobId)

	if err != nil {
		return fmt.Errorf("deleteStartTimes: %v", err)
	}
	return nil
}

func deleteStartMins(tx *sql.Tx, jobId int64) error {
	qry := "Delete From sched_window_run Where job_id=?"
	_, err := tx.Exec(qry, jobId)

	if err != nil {
		return fmt.Errorf("deleteStartTimes: %v", err)
	}
	return nil
}

func deleteRunWindow(tx *sql.Tx, jobId int64) error {
	if err := updateRunWindow(tx, jobId, "", ""); err != nil {
		return fmt.Errorf("deleteRunWindow: %w", err)
	}
	return nil
}

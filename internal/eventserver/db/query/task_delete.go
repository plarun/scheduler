package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	er "github.com/plarun/scheduler/pkg/error"
)

// DeleteTask deletes an existing task
func DeleteTask(tx *sql.Tx, tsk *task.TaskEntity) error {
	taskId, err := GetTaskId(tsk.Name())
	if err != nil {
		return err
	}

	// remove task run times
	if err := deleteStartTimes(tx, taskId); err != nil {
		return fmt.Errorf("DeleteTask: %w", err)
	}
	if err := deleteStartMins(tx, taskId); err != nil {
		return fmt.Errorf("DeleteTask: %w", err)
	}

	// remove all the relations of task
	if err := deleteTaskRelation(tx, taskId); err != nil {
		return fmt.Errorf("DeleteTask: %w", err)
	}

	// remove all the run history of task
	if err := ClearTaskRunHistory(tx, taskId); err != nil {
		return fmt.Errorf("DeleteTask: %w", err)
	}

	// remove the definition of task
	qry := "Delete From sched_task Where name=?"
	if _, err = tx.Exec(qry, tsk.Name); err != nil {
		return fmt.Errorf("DeleteTask: %w", er.NewDatabaseError(err.Error()))
	}
	return nil
}

func deleteStartTimes(tx *sql.Tx, taskId int64) error {
	qry := "Delete From sched_batch_run Where task_id=?"
	if _, err := tx.Exec(qry, taskId); err != nil {
		return fmt.Errorf("deleteStartTimes: %w", er.NewDatabaseError(err.Error()))
	}
	return nil
}

func deleteStartMins(tx *sql.Tx, taskId int64) error {
	qry := "Delete From sched_window_run Where task_id=?"
	if _, err := tx.Exec(qry, taskId); err != nil {
		return fmt.Errorf("deleteStartTimes: %w", er.NewDatabaseError(err.Error()))
	}
	return nil
}

func deleteRunWindow(tx *sql.Tx, taskId int64) error {
	if err := updateRunWindow(tx, taskId, "", ""); err != nil {
		return fmt.Errorf("deleteRunWindow: %w", err)
	}
	return nil
}

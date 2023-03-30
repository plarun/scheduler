package query

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/internal/eventserver/db"
	er "github.com/plarun/scheduler/pkg/error"
)

// ClearTaskRunHistory removes all the run history of a given task id
func ClearTaskRunHistory(tx *sql.Tx, id int64) error {
	qry := "Delete From task_run_history Where task_id=?"

	if _, err := tx.Exec(qry, id); err != nil {
		return fmt.Errorf("ClearTaskRunHistory: %w", er.NewDatabaseError(err.Error()))
	}
	return nil
}

func AddTaskRun(id int64) error {
	db := db.GetDatabase()

	qry := `Insert Into sched_run_history (task_id, seq_id, start_time, end_time, status)
		Select id, 
			(
				Select Ifnull(Max(seq_id), 0) + 1
				From sched_run_history
				Where task_id=?
			) seq_id, 
			last_start_time, 
			last_end_time, 
			current_status
		From sched_task
		Where id=?`

	if r, err := db.Exec(qry, id, id); err != nil {
		return fmt.Errorf("AddTaskRun: failed to add run entry: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("AddTaskRun: run entry added to history for task id - %d", id)
	}
	return nil
}

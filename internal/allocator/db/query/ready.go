package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db"
)

func InsertReadyTask(id int64) error {
	db := db.GetDatabase()

	qry := `Insert Into sched_ready (
		task_id,
		sys_entry_date,
		priority
	)
	Select task_id, now(), priority
	From sched_queue
	Where task_id=? And lock_flag=1`

	r, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("InsertReadyTask: failed to move task from queue to ready queue: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("InsertReadyTask: %d - task id inserted into sched_ready", id)
	}

	if err := setTaskstatus(id, task.StateReady); err != nil {
		return fmt.Errorf("InsertReadyTask: %w", err)
	}
	return nil
}

func DeleteReadyTask(id int64) error {
	db := db.GetDatabase()

	qry := `Delete From sched_ready Where task_id=?`

	r, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("DeleteReadyTask: failed to delete task from ready: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("DeleteReadyTask: %d - task id removed from sched_ready", id)
	}
	return nil
}

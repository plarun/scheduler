package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db"
	er "github.com/plarun/scheduler/pkg/error"
)

func InsertWaitTask(id int64) error {
	db := db.GetDatabase()

	qry := `Insert Into sched_wait (
			task_id,
			sys_entry_date,
			priority
		)
		Select task_id, now(), priority
		From sched_queue
		Where task_id=? And lock_flag=1`

	result, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("InsertWaitTask: failed to move task from queue to wait: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := result.RowsAffected(); n > 0 {
		log.Printf("InsertWaitTask: %d - task id inserted into sched_wait", id)
	}

	if err := setTaskstatus(id, task.StateWaiting); err != nil {
		return fmt.Errorf("InsertReadyTask: %w", err)
	}
	return nil
}

func DeleteWaitTask(id int64) error {
	db := db.GetDatabase()

	qry := `Delete From sched_wait Where task_id=?`

	result, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("DeleteWaitTask: failed to delete task from wait: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := result.RowsAffected(); n > 0 {
		log.Printf("DeleteWaitTask: %d - task id removed from sched_wait", id)
	}
	return nil
}

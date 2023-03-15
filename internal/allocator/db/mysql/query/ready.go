package query

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func InsertReadyTask(id int) error {
	db := mysql.GetDatabase()

	qry := `Insert Into sched_ready (
		task_id,
		sys_entry_date,
		priority
	)
	Select task_id, now(), priority
	From sched_queue
	Where task_id=? And lock_flag=1`

	_, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("InsertReadyTask: failed to move task from queue to ready queue: %v", err)
	}

	if err := setTaskstatus(id, task.StateReady); err != nil {
		return fmt.Errorf("InsertReadyTask: %w", err)
	}
	return nil
}

func DeleteReadyTask(id int) error {
	db := mysql.GetDatabase()

	qry := `Delete From sched_ready Where task_id=?`

	_, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("DeleteReadyTask: failed to delete task from ready: %v", err)
	}
	return nil
}

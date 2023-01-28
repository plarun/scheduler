package query

import (
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func InsertWaitTask(id int) error {
	db := mysql.GetDatabase()

	qry := `Insert Into sched_wait (
			task_id,
			sys_entry_date,
			priority
		)
		Select task_id, now(), priority
		From sched_queue
		Where task_id=? And flag=1`

	_, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("InsertWaitTask: failed to move task from queue to wait: %v", err)
	}
	return nil
}

func DeleteWaitTask(id int) error {
	db := mysql.GetDatabase()

	qry := `Delete From sched_wait Where task_id=?`

	_, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("DeleteWaitTask: failed to delete task from wait: %v", err)
	}
	return nil
}

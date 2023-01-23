package query

import (
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func InsertWaitJob(jobId int) error {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Insert Into sched_wait (
			job_id,
			sys_entry_date,
			priority
		)
		Select job_id, now(), priority
		From sched_queue
		Where job_id=? And flag=1`

	_, err := db.DB.Exec(qry, jobId)
	if err != nil {
		return fmt.Errorf("DequeueJob: failed to move job from queue to wait: %v", err)
	}
	return nil
}

func DeleteWaitJob(jobId int) error {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Delete From sched_wait Where job_id=?`

	_, err := db.DB.Exec(qry, jobId)
	if err != nil {
		return fmt.Errorf("DeleteWaitJob: failed to delete job from wait: %v", err)
	}
	return nil
}

package query

import (
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func InsertReadyjob(jobId int) error {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Insert Into sched_ready (
		job_id,
		sys_entry_date,
		priority
	)
	Select job_id, now(), priority
	From sched_queue
	Where job_id=? And flag=1`

	_, err := db.DB.Exec(qry, jobId)
	if err != nil {
		return fmt.Errorf("InsertReadyjob: failed to move job from queue to ready queue: %v", err)
	}
	return nil
}

func DeleteReadyJob(jobId int) error {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Delete From sched_ready Where job_id=?`

	_, err := db.DB.Exec(qry, jobId)
	if err != nil {
		return fmt.Errorf("DeleteReadyJob: failed to delete job from ready: %v", err)
	}
	return nil
}

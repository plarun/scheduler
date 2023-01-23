package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func QueueJobs() (int64, error) {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Insert Into sched_queue (
			job_id,
			sys_entry_date,
			priority
		)
		Select job_id
		From sched_stage
		Where is_bundle=0 And flag=1
		Union
		Select j.job_id 
		From sched_job j
			Inner Join sched_stage s On (j.parent_id=s.job_id)
		Where s.flag=1 And s.is_bundle=1`

	result, err := db.DB.Exec(qry)
	if err != nil {
		return 0, fmt.Errorf("QueueJobs: failed to push jobs into queue: %v", err)
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("QueueJobs: %v", err)
	}

	log.Printf("%d jobs pushed into queue", cnt)
	return cnt, nil
}

func DequeueJob(jobId int) error {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Delete From sched_queue Where job_id=?`

	_, err := db.DB.Exec(qry, jobId)
	if err != nil {
		return fmt.Errorf("DequeueJob: failed to delete job from queue: %v", err)
	}
	return nil
}

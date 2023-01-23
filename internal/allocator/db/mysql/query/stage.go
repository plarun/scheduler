package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func InsertStageJobs(window int) (int64, error) {
	db := mysql.GetDatabase()

	qry := `Insert Into sched_stage (
			job_id,
			sys_entry_date,
			priority,
			flag,
			is_bundle
		)
		Select
			j.job_id, 
			now(), 
			j.priority, 
			0, 
			Case When j.job_type = 'BUNDLE' Then 1 Else 0 End Case
		From sched_job j 
			Inner Join sched_batch_run b On (j.job_id=b.job_id)
		Where .run_flag=2 
			And b.start_time Between current_time And timestampadd(Second, ?, current_time)
		Union
		Select
			j.job_id, 
			now(), 
			j.priority, 
			0, 
			Case When j.job_type = 'BUNDLE' Then 1 Else 0 End Case
		From sched_job j 
			Inner Join sched_window_run w On (j.job_id=w.job_id)
		Where j.run_flag=3 
			And current_time Between j.start_window And j.end_window
			And start_min = minute(current_time)`

	db.Lock()
	defer db.Unlock()

	result, err := db.DB.Exec(qry, window)
	if err != nil {
		return 0, fmt.Errorf("PollJobs: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("PollJobs: %v", err)
	}

	log.Printf("%d jobs polled", cnt)
	return cnt, nil
}

func UpdateStageFlag(from, to int) (int64, error) {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Update sched_stage Set flag=? Where flag=?`

	result, err := db.DB.Exec(qry, to, from)
	if err != nil {
		return 0, fmt.Errorf("UpdateStageFlag: failed to update flags: %v", err)
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("QueueJobs: %v", err)
	}

	log.Printf("flag from %d to %d updated for %d staged jobs", from, to, cnt)
	return cnt, nil
}

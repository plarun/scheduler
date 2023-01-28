package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func StageTasks(window int) (int64, error) {
	db := mysql.GetDatabase()

	qry := `Insert Into sched_stage (
			task_id,
			sys_entry_date,
			priority,
			flag,
			is_bundle
		)
		Select
			t.id, now(), t.priority, 0, 
			Case When t.type = 'bundle' Then 1 Else 0 End
		From sched_task t
			Inner Join sched_batch_run b On (t.id=b.task_id)
		Where t.run_flag='batch'
			And b.start_time Between current_time And timestampadd(Second, ?, current_time)
		Union
		Select
			t.id, now(), t.priority, 0, 
			Case When t.type = 'bundle' Then 1 Else 0 End
		From sched_task t
			Inner Join sched_window_run w On (t.id=w.task_id)
		Where t.run_flag='window'
			And current_time Between t.start_window And t.end_window
			And start_min = minute(current_time)`

	result, err := db.DB.Exec(qry, window)
	if err != nil {
		return 0, fmt.Errorf("StageTasks: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("StageTasks: %v", err)
	}

	log.Printf("%d tasks staged", cnt)
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
		return 0, fmt.Errorf("UpdateStageFlag: %v", err)
	}

	log.Printf("flag from %d to %d updated for %d staged tasks", from, to, cnt)
	return cnt, nil
}

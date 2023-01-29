package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

var (
	stableStates string = fmt.Sprintf("%s,%s,%s,%s", task.StateIdle, task.StateAborted, task.StateFailure, task.StateSuccess)
)

// StageTasks inserts the tasks which are scheduled for given cycle
// into 'sched_stage'. When staging the tasks, flag will be 0 which
// indicates newly staged task.
func StageTasks(window int) error {
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
			And b.current_status In (` + stableStates + `)
		Union
		Select
			t.id, now(), t.priority, 0,
			Case When t.type = 'bundle' Then 1 Else 0 End
		From sched_task t
			Inner Join sched_window_run w On (t.id=w.task_id)
		Where t.run_flag='window'
			And current_time Between t.start_window And t.end_window
			And start_min = minute(current_time)
			And b.current_status In (` + stableStates + `)`

	result, err := db.DB.Exec(qry, window)
	if err != nil {
		return fmt.Errorf("StageTasks: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("StageTasks: %v", err)
	}

	log.Printf("%d tasks staged", cnt)
	return nil
}

// SetStagedStatus changes the status of newly staged tasks to 'staged'
func SetStagedStatus() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t Join sched_stage s On t.id=s.task_id
	Set t.current_status='?'
	Where s.flag=0`

	result, err := db.DB.Exec(qry, task.StateStaged)
	if err != nil {
		return fmt.Errorf("SetStagedStatus: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SetStagedStatus: %v", err)
	}

	log.Printf("%d tasks changed to status staged", cnt)
	return nil
}

// MigrateStage completes the staging and sets the tasks to ready for queuing
func MigrateStage() error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage s Join sched_task t On s.task_id=t.id
	Set s.flag=1
	Where s.flag=0 And t.current_status='?'`

	result, err := db.DB.Exec(qry, task.StateStaged)
	if err != nil {
		return fmt.Errorf("MigrateStage: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("MigrateStage: %v", err)
	}

	log.Printf("%d tasks ready for migration to queue", cnt)
	return nil
}

// UpdateStageFlag changes the flag of the staged task
func UpdateStageFlag(from, to int) (int64, error) {
	db := mysql.GetDatabase()

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

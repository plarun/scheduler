package query

import (
	"fmt"
	"log"
	"time"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// LockForStaging locks the task for staging
func LockForStaging(cycle time.Duration) error {
	db := mysql.GetDatabase()

	sec := int(cycle / time.Second)

	qry := `With tasks As (
		Select t.id
		From sched_task t
			Inner Join sched_batch_run b On (t.id=b.task_id)
		Where t.run_flag='batch'
			And date_format(b.start_time,'%H:%i') >= date_format(now(),'%H:%i')
			And date_format(b.start_time,'%H:%i') < date_format(timestampadd(Minute, 1, now()), '%H:%i')
			And t.current_status In ('idle','success','failure','aborted')
			And t.lock_flag=0
			And (
				t.last_end_time Is Null
				Or t.last_end_time < date_format(now(),'%H:%i'))
		Union All
		Select t.id
		From sched_task t
			Inner Join sched_window_run w On (t.id=w.task_id)
		Where t.run_flag='window'
			And current_time Between t.start_window And t.end_window
			And start_min = minute(current_time)
			And t.current_status In ('idle','success','failure','aborted')
			And t.lock_flag=0
			And (
				t.last_end_time Is Null
				Or t.last_end_time < date_format(now(),'%H:%i'))
		) Update sched_task
		Set lock_flag=1
		Where id In (Select id From tasks)`

	result, err := db.DB.Exec(qry, sec)
	if err != nil {
		return fmt.Errorf("LockForStaging: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("LockForStaging: %v", err)
	}

	log.Printf("%d tasks locked for staging", cnt)
	return nil
}

// StageTasks inserts the tasks which are scheduled for given cycle
// into 'sched_stage'. When staging the tasks, flag will be 0 which
// indicates newly staged task.
func StageTasks() error {
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
		Where lock_flag=?
			And current_status<>?`

	result, err := db.DB.Exec(qry, 1, string(task.StateQueued))
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
	Set t.current_status=?
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

// SetStagedFlag completes the staging and sets the tasks to ready for queuing
func SetStagedFlag() error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage s Join sched_task t On s.task_id=t.id
	Set s.flag=1
	Where s.flag=0 And t.current_status=?`

	result, err := db.DB.Exec(qry, task.StateStaged)
	if err != nil {
		return fmt.Errorf("SetStagedFlag: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SetStagedFlag: %v", err)
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

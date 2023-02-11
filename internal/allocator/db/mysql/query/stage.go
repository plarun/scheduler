package query

import (
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// LockForStaging locks the tasks for staging which are having scheduled run
// either batch run or window run
func LockForStaging() error {
	db := mysql.GetDatabase()

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

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("LockForStaging: %w", err)
	}
	return nil
}

// StageLockedTasks inserts the tasks which are scheduled for given cycle
// into 'sched_stage'. When staging the tasks, flag will be 0 which
// indicates newly staged task.
func StageLockedTasks() error {
	db := mysql.GetDatabase()

	qry := `Insert Into sched_stage (
			task_id,
			sys_entry_date,
			priority,
			flag,
			is_bundle
		)
		Select
			id, now(), priority, 0, Case When type = 'bundle' Then 1 Else 0 End
		From sched_task
		Where lock_flag=?
			And current_status Not In ('staged', 'queued', 'ready', 'waiting', 'running')`

	if _, err := db.DB.Exec(qry, 1); err != nil {
		return fmt.Errorf("StageLockedTasks: %w", err)
	}
	return nil
}

// MarkAsStaged changes the status of newly staged tasks to 'staged'
func MarkAsStaged() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t 
	Join sched_stage s On t.id=s.task_id
	Set t.current_status='staged'
	Where s.flag=0`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("MarkAsStaged: %w", err)
	}
	return nil
}

// SetStagedFlag completes the staging and sets the tasks to ready for queuing
func SetStagedFlag() error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage s Join sched_task t On s.task_id=t.id
	Set s.flag=1
	Where s.flag=0 And t.current_status='staged'`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("SetStagedFlag: %w", err)
	}
	return nil
}

// LockStagedBundles locks the staged bundle tasks for staging its tasks
func LockStagedBundles() error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage
		Set flag=4
		Where flag=2 
			And is_bundle=1`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("LockStagedBundles: %v", err)
	}
	return nil
}

// LockBundledTasksForStaging locks the bundled tasks for staging
// whose bundles are already staged
func LockBundledTasksForStaging() error {
	db := mysql.GetDatabase()

	qry := `With tasks As (
			Select t.id
			From sched_task t
				Inner Join sched_stage s On (t.parent_id=s.task_id)
			Where s.is_bundle=1 And s.flag=4
		) Update sched_task
		Set lock_flag=1
		Where id In (Select id From tasks)`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("StageBundledTasks: failed to stage the tasks under bundle: %v", err)
	}
	return nil
}

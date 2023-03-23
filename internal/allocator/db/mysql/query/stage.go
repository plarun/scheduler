package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
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
				Or Timestampdiff(Minute, t.last_end_time, now()))
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
				Or Timestampdiff(Minute, t.last_end_time, now()))
		) Update sched_task
		Set lock_flag=1
		Where id In (Select id From tasks)`

	if r, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("LockForStaging: %w", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("LockForStaging: %d tasks are locked for scheduling", n)
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

	if r, err := db.DB.Exec(qry, 1); err != nil {
		return fmt.Errorf("StageLockedTasks: %w", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("StageLockedTasks: %d tasks are locked after staging", n)
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

	if r, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("MarkAsStaged: %w", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("MarkAsStaged: %d tasks are marked as staged", n)
	}
	return nil
}

// SetStagedFlag completes the staging and sets the tasks to ready for queuing
func SetStagedFlag() error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage s Join sched_task t On s.task_id=t.id
	Set s.flag=1
	Where s.flag=0 And t.current_status='staged'`

	if r, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("SetStagedFlag: %w", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetStagedFlag: %d tasks are flaged as staged", n)
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

	if r, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("LockBundledTasksForStaging: failed to stage the tasks under bundle: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("LockBundledTasksForStaging: %d tasks of bundle are locked as staged", n)
	}
	return nil
}

// ChangeStagedBundleLock changes the lock of staged bundles
func ChangeStagedBundleLock(from, to int) error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage
		Set flag=?
		Where flag=?
			And is_bundle=1`

	if r, err := db.DB.Exec(qry, to, from); err != nil {
		return fmt.Errorf("ChangeStagedBundleLock: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("ChangeStagedBundleLock: %d bundle tasks are changed from flag %d to %d", n, from, to)
	}
	return nil
}

func MarkBundleAsRunning() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task
		Set current_status=?, last_start_time=now(), last_end_time=null
		Where id In (
			Select task_id
			From sched_stage
			Where flag=5 And is_bundle=1)`

	if r, err := db.DB.Exec(qry, string(task.StateRunning)); err != nil {
		return fmt.Errorf("MarkBundleAsRunning: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("MarkBundleAsRunning: %d bundle tasks status set to running", n)
	}
	return nil
}

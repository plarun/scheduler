package query

import (
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// LockForEnqueue locks the staged task for queuing
// so it will be considered for moving into queue
func LockForEnqueue() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t Join sched_stage s On t.id=s.task_id
		Set s.flag=2
		Where s.flag=1
			And t.current_status='staged'`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("LockForEnqueue: %w", err)
	}

	return nil
}

// EnqueueTasks inserts the staged tasks into sched_queue.
func EnqueueTasks() error {
	db := mysql.GetDatabase()

	qry := `Insert Into sched_queue (
			task_id,
			sys_entry_date,
			priority
		)
		Select task_id, now(), priority
		From sched_stage
		Where is_bundle=0 And flag=2`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("EnqueueTasks: failed to push tasks into queue: %v", err)
	}
	return nil
}

// SetQueueStatus sets the state of queued tasks to queued
func SetQueueStatus() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t Join sched_queue q On t.id=q.task_id
		Set t.current_status='queued'
		Where t.current_status='staged'`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("SetQueueStatus: %w", err)
	}

	return nil
}

// SetQueuedFlag changes the flag of task in stage to Queued
func SetQueuedFlag() error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage s Join sched_task t On s.task_id=t.id
	Set s.flag=3
	Where s.flag=2 And t.current_status='queued'`

	if _, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("SetQueuedFlag: %v", err)
	}

	return nil
}

// LockForDequeue locks a task in queue for dequeue
func LockForDequeue(id int) error {
	db := mysql.GetDatabase()

	qry := `Update sched_queue
	Set lock_flag=1
	Where task_id=?`

	if _, err := db.DB.Exec(qry, id); err != nil {
		return fmt.Errorf("LockForDequeue: %v", err)
	}
	return nil
}

// RemoveTask removes a task from queue
func RemoveTask(id int) error {
	db := mysql.GetDatabase()

	qry := `Delete From sched_queue Where task_id=?`

	if _, err := db.DB.Exec(qry, id); err != nil {
		return fmt.Errorf("RemoveTask: failed to delete task from queue: %v", err)
	}
	return nil
}

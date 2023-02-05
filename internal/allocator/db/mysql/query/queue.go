package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// LockForEnqueue changes the flag of staged tasks to pre queue state
func LockForEnqueue() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t Join sched_stage s On t.id=s.task_id
		Set s.flag=2
		Where s.flag=1
			And t.current_status='staged'`

	result, err := db.DB.Exec(qry)
	if err != nil {
		return fmt.Errorf("LockForEnqueue: %w", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("LockForEnqueue: %w", err)
	}

	log.Printf("%d staged tasks are locked for queuing", cnt)
	return nil
}

// EnqueueTasks inserts the staged tasks into sched_stage.
func EnqueueTasks() error {
	db := mysql.GetDatabase()

	qry := `Insert Into sched_queue (
			task_id,
			sys_entry_date,
			priority
		)
		Select task_id
		From sched_stage
		Where is_bundle=0 And flag=2
		Union All
		Select t.task_id 
		From sched_task t
			Inner Join sched_stage s On (t.parent_id=s.task_id)
		Where s.is_bundle=1 And s.flag=2`

	result, err := db.DB.Exec(qry)
	if err != nil {
		return fmt.Errorf("EnqueueTasks: failed to push tasks into queue: %v", err)
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("EnqueueTasks: %v", err)
	}

	log.Printf("%d tasks pushed into queue", cnt)
	return nil
}

// SetQueueStatus sets the state of queued tasks to queued
func SetQueueStatus() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t Join sched_queue q On t.id=q.task_id
		Set t.current_status='queued'
		Where t.current_status='staged'`

	_, err := db.DB.Exec(qry)
	if err != nil {
		return fmt.Errorf("SetQueueStatus: %w", err)
	}

	return nil
}

// SetQueuedFlag changes the flag of task in stage to Queued
func SetQueuedFlag() error {
	db := mysql.GetDatabase()

	qry := `Update sched_queue s Join sched_task t On s.task_id=t.id
	Set s.flag=3
	Where s.flag=2 And t.current_status='queued'`

	_, err := db.DB.Exec(qry)
	if err != nil {
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

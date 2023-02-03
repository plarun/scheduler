package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// SetFlagPreQueue changes the flag of staged tasks to pre queue state
func SetFlagPreQueue() error {
	cnt, err := UpdateStageFlag(1, 2)
	if err != nil {
		return fmt.Errorf("MigrateStage: %v", err)
	}

	log.Printf("%d tasks ready for migration to queue", cnt)
	return nil
}

// SetStatusQueue changes the task status to 'queued'
func SetStatusQueue() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t Join sched_stage s On t.id=s.task_id
	Set t.current_status='?'
	Where s.flag=2`

	result, err := db.DB.Exec(qry, task.StateQueued)
	if err != nil {
		return fmt.Errorf("SetStatusQueue: %v", err)
	}

	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("SetStatusQueue: %v", err)
	}

	log.Printf("%d tasks changed to status queued", cnt)
	return nil
}

// QueueTasks inserts the staged tasks into sched_stage.
func QueueTasks() error {
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
		return fmt.Errorf("QueueTasks: failed to push tasks into queue: %v", err)
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("QueueTasks: %v", err)
	}

	log.Printf("%d tasks pushed into queue", cnt)
	return nil
}

// SetFlagPostQueue changes the flag of task in stage to Queued
func SetFlagPostQueue() error {
	cnt, err := UpdateStageFlag(2, 3)
	if err != nil {
		return fmt.Errorf("SetFlagPostQueue: %v", err)
	}

	log.Printf("%d tasks migrated to queue", cnt)
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

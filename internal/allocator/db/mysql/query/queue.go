package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

const (
	QueueLockNew      int = 0
	QueueLockChecking int = 1
	QueueLockReady    int = 2
	QueueLockWait     int = 3
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

// LockForConditionCheck locks the queued tasks in sched_queue for start condition check
func LockForConditionCheck() error {
	db := mysql.GetDatabase()

	qry := `With tasks As (
			Select q.task_id
			From sched_task t, sched_queue q
			Where t.id=q.task_id
				And q.lock_flag=0
				And t.current_status='queued'
		) Update sched_queue
		Set lock_flag=?
		Where task_id In (
			Select task_id 
			From tasks)`

	if _, err := db.DB.Exec(qry, QueueLockChecking); err != nil {
		return fmt.Errorf("LockForConditionCheck: %v", err)
	}
	return nil
}

func PickQueueLockedTasks() ([]int, error) {
	db := mysql.GetDatabase()

	qry := `Select task_id
		From sched_queue
		Where lock_flag=?`

	res := make([]int, 0)

	rows, err := db.DB.Query(qry, QueueLockChecking)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("PickQueueLockedTasks: %w", err)
	}

	for rows.Next() {
		var taskId int
		rows.Scan(&taskId)
		res = append(res, taskId)
	}
	return res, nil
}

// SetQueueLockFlag sets the given lock flag on queued task
func SetQueueLockFlag(id, flag int) error {
	db := mysql.GetDatabase()

	qry := `Update sched_queue
	Set lock_flag=?
	Where task_id=?`

	if _, err := db.DB.Exec(qry, id, flag); err != nil {
		return fmt.Errorf("LockForDequeue: %v", err)
	}
	return nil
}

func MoveQueueToReady(id int) error {
	if err := InsertReadyTask(id); err != nil {
		return fmt.Errorf("MoveQueueToReady: %w", err)
	}

	if err := DequeueTask(id); err != nil {
		return fmt.Errorf("MoveQueueToReady: %w", err)
	}
	return nil
}

func MoveQueueToWait(id int) error {
	if err := InsertWaitTask(id); err != nil {
		return fmt.Errorf("MoveQueueToWait: %w", err)
	}

	if err := DequeueTask(id); err != nil {
		return fmt.Errorf("MoveQueueToWait: %w", err)
	}
	return nil
}

// DequeueTask removes a task from queue
func DequeueTask(id int) error {
	db := mysql.GetDatabase()

	qry := `Delete From sched_queue Where task_id=?`

	if _, err := db.DB.Exec(qry, id); err != nil {
		return fmt.Errorf("DequeueTask: failed to remove task from queue: %v", err)
	}
	return nil
}

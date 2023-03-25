package query

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
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

	if r, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("LockForEnqueue: %w", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("LockForEnqueue: %d tasks are locked for queueing", n)
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

	if r, err := db.DB.Exec(qry); err != nil {
		return fmt.Errorf("EnqueueTasks: failed to push tasks into queue: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("EnqueueTasks: %d tasks are queued into sched_queue", n)
	}
	return nil
}

// SetQueueStatus sets the state of queued tasks to queued
func SetQueueStatus() error {
	db := mysql.GetDatabase()

	qry := `Update sched_task t Join sched_queue q On t.id=q.task_id
		Set t.current_status=?
		Where t.current_status=?`

	if r, err := db.DB.Exec(qry, string(task.StateQueued), string(task.StateStaged)); err != nil {
		return fmt.Errorf("SetQueueStatus: %w", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetQueueStatus: %d tasks are set to status queued", n)
	}

	return nil
}

// SetQueuedFlag changes the flag of task in stage to Queued
func SetQueuedFlag() error {
	db := mysql.GetDatabase()

	qry := `Update sched_stage s Join sched_task t On s.task_id=t.id
		Set s.flag=3
		Where s.flag=2 And t.current_status=?`

	if r, err := db.DB.Exec(qry, string(task.StateQueued)); err != nil {
		return fmt.Errorf("SetQueuedFlag: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetQueuedFlag: %d tasks are flaged as queued", n)
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
				And t.current_status=?
		) Update sched_queue
		Set lock_flag=?
		Where task_id In (
			Select task_id 
			From tasks)`

	if r, err := db.DB.Exec(qry, string(task.StateQueued), QueueLockChecking); err != nil {
		return fmt.Errorf("LockForConditionCheck: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("LockForConditionCheck: %d tasks are locked for condition check", n)
	}
	return nil
}

func PickQueueLockedTasks() ([]int64, error) {
	db := mysql.GetDatabase()

	qry := `Select task_id
		From sched_queue
		Where lock_flag=?`

	res := make([]int64, 0)

	rows, err := db.DB.Query(qry, QueueLockChecking)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("PickQueueLockedTasks: %w", err)
	}

	for rows.Next() {
		var taskId int64
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

	if r, err := db.DB.Exec(qry, id, flag); err != nil {
		return fmt.Errorf("SetQueueLockFlag: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetQueueLockFlag: %d - task id in sched_queue set to flag %d", id, flag)
	}
	return nil
}

func MoveQueueToReady(id int64) error {
	if err := InsertReadyTask(id); err != nil {
		return fmt.Errorf("MoveQueueToReady: %w", err)
	}

	if err := DequeueTask(id); err != nil {
		return fmt.Errorf("MoveQueueToReady: %w", err)
	}
	return nil
}

func MoveQueueToWait(id int64) error {
	if err := InsertWaitTask(id); err != nil {
		return fmt.Errorf("MoveQueueToWait: %w", err)
	}

	if err := DequeueTask(id); err != nil {
		return fmt.Errorf("MoveQueueToWait: %w", err)
	}
	return nil
}

// DequeueTask removes a task from queue
func DequeueTask(id int64) error {
	db := mysql.GetDatabase()

	qry := `Delete From sched_queue Where task_id=?`

	if r, err := db.DB.Exec(qry, id); err != nil {
		return fmt.Errorf("DequeueTask: failed to remove task from queue: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("DequeueTask: %d - task id removed from sched_queue", id)
	}
	return nil
}

// DependentConditionCheck locks the queued tasks in sched_queue for start condition check
func MoveDependentToQueue(ctx context.Context, id int64) error {
	tx, err := mysql.GetDatabase().DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("MoveDependentToQueue: %v", err)
	}

	if err := moveWaitToQueue(tx, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("MoveDependentToQueue: %w", err)
	}

	if err := setStatusAfterAwaken(tx, id, task.StateQueued); err != nil {
		tx.Rollback()
		return fmt.Errorf("MoveDependentToQueue: %w", err)
	}

	if err := clearWaitAfterAwaken(tx, id); err != nil {
		tx.Rollback()
		return fmt.Errorf("MoveDependentToQueue: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("MoveDependentToQueue: %v", err)
	}

	return nil
}

// move dependent tasks from waiting area to queue
func moveWaitToQueue(tx *sql.Tx, id int64) error {
	qry := `Insert Into sched_queue (task_id, sys_entry_date, priority)
		Select task_id, now(), priority
		From sched_wait 
		Where task_id in (
			Select task_id
			From sched_task_relation
			Where cond_task_id=?)`

	if r, err := tx.Exec(qry, id); err != nil {
		return fmt.Errorf("moveWaitToQueue: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("moveWaitToQueue: %d tasks are moved to queue from wait", n)
	}
	return nil
}

// setStatusAfterAwaken set status to queued after awaken
func setStatusAfterAwaken(tx *sql.Tx, id int64, state task.State) error {
	qry := `Update sched_task
		Set current_status=? 
		Where id in (
			Select task_id
			From sched_task_relation
			Where cond_task_id=?)`

	if r, err := tx.Exec(qry, string(state), id); err != nil {
		return fmt.Errorf("setStatusAfterAwaken: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("setStatusAfterAwaken: %d tasks are set to status queued", n)
	}
	return nil
}

func clearWaitAfterAwaken(tx *sql.Tx, id int64) error {
	qry := `Delete From sched_wait
	Where task_id In (
		Select task_id
		From sched_task_relation
		Where cond_task_id=?
	)`

	if r, err := tx.Exec(qry, id); err != nil {
		return fmt.Errorf("clearWaitAfterAwaken: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("clearWaitAfterAwaken: %d tasks are cleared from wait", n)
	}
	return nil
}

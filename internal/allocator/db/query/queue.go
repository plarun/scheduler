package query

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db"
	er "github.com/plarun/scheduler/pkg/error"
)

const (
	QueueFlag_New         int = 0
	QueueFlag_Checking    int = 1
	QueueFlag_MoveToReady int = 2
	QueueFlag_MoveToWait  int = 3
)

// LockForEnqueue locks the staged task for queuing
// so it will be considered for moving into queue
func LockForEnqueue() error {
	db := db.GetDatabase()

	qry := `Update sched_task t Join sched_stage s On t.id=s.task_id
		Set s.flag=?
		Where s.flag=?
			And t.current_status='staged'`

	if r, err := db.DB.Exec(qry, StageFlag_CTaskQueuing, StageFlag_CTaskStaged); err != nil {
		return fmt.Errorf("LockForEnqueue: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("LockForEnqueue: %d tasks are locked for queueing", n)
	}

	return nil
}

// EnqueueTasks inserts the staged tasks into sched_queue.
func EnqueueTasks() error {
	db := db.GetDatabase()

	qry := `Insert Into sched_queue (
			task_id,
			sys_entry_date,
			priority
		)
		Select task_id, now(), priority
		From sched_stage
		Where is_bundle=0 And flag=?`

	if r, err := db.DB.Exec(qry, StageFlag_CTaskQueuing); err != nil {
		return fmt.Errorf("EnqueueTasks: failed to push tasks into queue: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("EnqueueTasks: %d tasks are queued into sched_queue", n)
	}
	return nil
}

// SetQueueStatus sets the state of queued tasks to queued
func SetQueueStatus() error {
	db := db.GetDatabase()

	qry := `Update sched_task t Join sched_queue q On t.id=q.task_id
		Set t.current_status=?
		Where t.current_status=?`

	if r, err := db.DB.Exec(qry, string(task.StateQueued), string(task.StateStaged)); err != nil {
		return fmt.Errorf("SetQueueStatus: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetQueueStatus: %d tasks are set to status queued", n)
	}

	return nil
}

// SetQueuedFlag changes the flag of task in stage to Queued
func SetQueuedFlag() error {
	db := db.GetDatabase()

	qry := `Update sched_stage s Join sched_task t On s.task_id=t.id
		Set s.flag=?
		Where s.flag=? And t.current_status=?`

	if r, err := db.DB.Exec(qry, StageFlag_CTaskQueued, StageFlag_CTaskQueuing, string(task.StateQueued)); err != nil {
		return fmt.Errorf("SetQueuedFlag: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetQueuedFlag: %d tasks are flaged as queued", n)
	}
	return nil
}

// LockForConditionCheck locks the queued tasks in sched_queue for start condition check
func LockForConditionCheck() error {
	db := db.GetDatabase()

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

	if r, err := db.DB.Exec(qry, string(task.StateQueued), QueueFlag_Checking); err != nil {
		return fmt.Errorf("LockForConditionCheck: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("LockForConditionCheck: %d tasks are locked for condition check", n)
	}
	return nil
}

func PickQueueLockedTasks() ([]int64, error) {
	db := db.GetDatabase()

	qry := `Select task_id
		From sched_queue
		Where lock_flag=?`

	res := make([]int64, 0)

	rows, err := db.DB.Query(qry, QueueFlag_Checking)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("PickQueueLockedTasks: %w", er.NewDatabaseError(err.Error()))
	}

	for rows.Next() {
		var taskId int64
		if err := rows.Scan(&taskId); err != nil {
			return res, fmt.Errorf("PickQueueLockedTasks: %w", er.NewDatabaseError(err.Error()))
		}
		res = append(res, taskId)
	}
	return res, nil
}

// SetQueueLockFlag sets the given lock flag on queued task.
// func SetQueueLockFlag(id, flag int) error {
// 	db := db.GetDatabase()

// 	qry := `Update sched_queue
// 	Set lock_flag=?
// 	Where task_id=?`

// 	if r, err := db.DB.Exec(qry, flag, id); err != nil {
// 		return fmt.Errorf("SetQueueLockFlag: %v", err)
// 	} else if n, _ := r.RowsAffected(); n > 0 {
// 		log.Printf("SetQueueLockFlag: %d - task id in sched_queue set to flag %d", id, flag)
// 	}
// 	return nil
// }

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
	db := db.GetDatabase()

	qry := `Delete From sched_queue Where task_id=?`

	if r, err := db.DB.Exec(qry, id); err != nil {
		return fmt.Errorf("DequeueTask: failed to remove task from queue: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("DequeueTask: %d - task id removed from sched_queue", id)
	}
	return nil
}

// DependentConditionCheck locks the queued tasks in sched_queue for start condition check
func MoveDependentToQueue(ctx context.Context, id int64) error {
	tx, err := db.GetDatabase().DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("MoveDependentToQueue: %w", er.NewDatabaseError(err.Error()))
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
		return fmt.Errorf("MoveDependentToQueue: %w", er.NewDatabaseError(err.Error()))
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
		return fmt.Errorf("moveWaitToQueue: %w", er.NewDatabaseError(err.Error()))
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
		return fmt.Errorf("setStatusAfterAwaken: %w", er.NewDatabaseError(err.Error()))
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
		return fmt.Errorf("clearWaitAfterAwaken: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("clearWaitAfterAwaken: %d tasks are cleared from wait", n)
	}
	return nil
}

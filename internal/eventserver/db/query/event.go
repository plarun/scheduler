package query

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/eventserver/db"
	er "github.com/plarun/scheduler/pkg/error"
)

func SetTaskStatusByEvent(id int64, event task.SendEvent) error {
	db := db.GetDatabase()

	// create a transaction
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("SetTaskStatusByEvent: %w", err)
	}

	if err := eventStateChange(tx, id, event); err != nil {
		tx.Rollback()
		return fmt.Errorf("SetTaskStatusByEvent: %w", err)
	} else {
		tx.Commit()
		return nil
	}
}

func eventStateChange(tx *sql.Tx, id int64, event task.SendEvent) error {
	if event.IsStart() {
		// stage
		if ok, err := addToStageByEvent(tx, id); err != nil {
			return fmt.Errorf("eventStateChange: %w", err)
		} else if !ok {
			return er.NewTaskEventError(er.ErrEventStartOnUnstableTask)
		}
		// ready
		if ok, err := addToReadyByEvent(tx, id); err != nil {
			return fmt.Errorf("eventStateChange: %w", err)
		} else if !ok {
			return er.NewTaskEventError(er.ErrEventStartOnUnstableTask)
		}
		// task status upd
		if err := SetTaskStatus(id, task.StateReady); err != nil {
			return fmt.Errorf("eventStateChange: %w", err)
		}
	} else if event.IsAbort() {
		// todo
		return nil
	} else if event.IsFreeze() || event.IsGreen() || event.IsRed() || event.IsReset() {
		var state string
		if event.IsGreen() {
			state = string(task.StateSuccess)
		} else if event.IsRed() {
			state = string(task.StateFailure)
		} else if event.IsReset() {
			state = string(task.StateIdle)
		} else if event.IsFreeze() {
			state = string(task.StateFrozen)
		}

		qry := `Update sched_task
		Set current_status=?, last_end_time=current_time
		Where id=? And lock_flag=0 And current_status Not In (?, ?, ?, ?, ?, ?, ?)`

		r, err := tx.Exec(qry,
			state,
			id,
			string(task.StateRunning),
			string(task.StateQueued),
			string(task.StateReady),
			string(task.StateStaged),
			string(task.StateWaiting))
		if err != nil {
			return fmt.Errorf("eventStateChange: %w", err)
		} else if n, _ := r.RowsAffected(); n > 0 {
			log.Printf("eventStateChange: %d - task id set to status %s", id, string(state))
		} else {
			if event.IsGreen() {
				return er.NewTaskEventError(er.ErrEventGreenOnUnstableTask)
			} else if event.IsRed() {
				return er.NewTaskEventError(er.ErrEventRedOnUnstableTask)
			} else if event.IsReset() {
				return er.NewTaskEventError(er.ErrEventResetOnUnstableTask)
			} else if event.IsFreeze() {
				return er.NewTaskEventError(er.ErrEventFreezeOnUnstableTask)
			}
			return nil
		}
	} else {
		// ignore
		return nil
	}
	return nil
}

func addToReadyByEvent(tx *sql.Tx, id int64) (bool, error) {
	qry := `Insert Into sched_ready (
		task_id,
		sys_entry_date,
		priority
	)
	Select id, now(), priority
	From sched_task
	Where id=?`

	r, err := tx.Exec(qry, id)
	if err != nil {
		return false, fmt.Errorf("addToReadyByEvent: failed to add task id=%d to sched_ready: %w", id, er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("addToReadyByEvent: task id=%d added to sched_ready", id)
		return true, nil
	} else {
		return true, nil
	}
}

func addToStageByEvent(tx *sql.Tx, id int64) (bool, error) {
	qry := `Insert Into sched_stage (
		task_id,
		sys_entry_date,
		priority,
		flag,
		is_bundle
	)
	Select
		id, now(), priority, 3, Case When type = 'bundle' Then 1 Else 0 End
	From sched_task
	Where id=?
		And current_status Not In ('staged', 'queued', 'ready', 'waiting', 'running')`

	if r, err := tx.Exec(qry, id); err != nil {
		return false, fmt.Errorf("AddToStage: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("AddToStage: task id=%d is staged", id)
		return true, nil
	} else {
		return true, nil
	}
}

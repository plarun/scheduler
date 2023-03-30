package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	er "github.com/plarun/scheduler/pkg/error"
)

// InsertTask inserts a new task
// task can be one of below category
// 1. runnable/non-runnable
// 2. bundle/callable
// 3. batch-run/window-run
func InsertTask(tx *sql.Tx, tsk *task.TaskEntity) error {
	// insert the task
	insertedTaskId, err := insertTask(tx, tsk)
	if err != nil {
		return fmt.Errorf("insertTask: %w", err)
	}

	// insert run times
	if f, ok := tsk.GetFieldStartTimes(); ok {
		// Start times of task
		if err := insertStartTimes(tx, insertedTaskId, f.Value()); err != nil {
			return fmt.Errorf("insertTask: %w", err)
		}
	} else if f, ok := tsk.GetFieldStartMins(); ok {
		// Start mins of task
		if err := insertStartMins(tx, insertedTaskId, f.Value()); err != nil {
			return fmt.Errorf("insertTask: %w", err)
		}
	}

	// insert task relation
	if _, ok := tsk.GetFieldCondition(); ok {
		f, _ := tsk.GetField(task.FIELD_CONDITION)
		distTasks := f.(*task.Condition).DistinctTasks()
		dependentTasksId, err := getTaskIdList(tx, distTasks)
		if err != nil {
			return fmt.Errorf("insertTask: %w", err)
		}

		if err := insertTaskRelation(tx, insertedTaskId, dependentTasksId); err != nil {
			return fmt.Errorf("insertTask: %w", err)
		}
	}
	return nil
}

// insertTask inserts a new task
func insertTask(tx *sql.Tx, tsk *task.TaskEntity) (int64, error) {
	var insertedTaskId int64
	var err error
	var result sql.Result

	var (
		command     sql.NullString = sql.NullString{Valid: false}
		condition   sql.NullString = sql.NullString{Valid: false}
		label       sql.NullString = sql.NullString{Valid: false}
		profile     sql.NullString = sql.NullString{Valid: false}
		outLogFile  sql.NullString = sql.NullString{Valid: false}
		errLogFile  sql.NullString = sql.NullString{Valid: false}
		startWindow sql.NullString = sql.NullString{Valid: false}
		endWindow   sql.NullString = sql.NullString{Valid: false}
		parent      sql.NullInt64  = sql.NullInt64{Valid: false}
		rundays     sql.NullInt32  = sql.NullInt32{Valid: false}

		priority int32
		tasktype string
	)

	if f, ok := tsk.GetFieldType(); ok {
		tasktype = string(f.Value())
	}

	if f, ok := tsk.GetFieldCommand(); ok {
		command.String, command.Valid = f.Value(), ok
	}

	if f, ok := tsk.GetFieldCondition(); ok {
		condition.String, condition.Valid = f.Value(), ok
	}

	if f, ok := tsk.GetFieldLabel(); ok {
		label.String, label.Valid = f.Value(), ok
	}

	if f, ok := tsk.GetFieldOutLogFile(); ok {
		outLogFile.String, outLogFile.Valid = f.Value(), ok
	}

	if f, ok := tsk.GetFieldErrLogFile(); ok {
		errLogFile.String, errLogFile.Valid = f.Value(), ok
	}

	if f, ok := tsk.GetFieldProfile(); ok {
		profile.String, profile.Valid = f.Value(), ok
	}

	if f, ok := tsk.GetFieldPriority(); ok {
		priority = f.Value()
	}

	if f, ok := tsk.GetFieldParent(); ok {
		id, err := GetTaskId(f.Value())
		if err != nil {
			return insertedTaskId, fmt.Errorf("insertTask: %w", err)
		}
		parent.Int64, parent.Valid = id, ok
	}

	if f, ok := tsk.GetFieldRunDays(); ok {
		rundays.Int32, rundays.Valid = int32(f.Value()), ok
	}

	runFlag := task.GetRunFlag(tsk)
	if runFlag.IsWindow() {
		if f, ok := tsk.GetFieldRunWindow(); ok {
			val := f.Value()
			startWindow.String, startWindow.Valid = val[0], ok
			endWindow.String, endWindow.Valid = val[1], ok
		}
	}

	qry := `Insert Into sched_task (
		parent_id,
		name,
		type,
		run_flag,
		start_condition,
		command,
		std_out_log,
		std_err_log,
		label,
		profile,
		run_days_bit,
		start_window,
		end_window,
		priority,
		created_on,
		created_by,
		current_status
		) Values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,now(),?,?)`

	result, err = tx.Exec(qry,
		parent,
		tsk.Name(),
		tasktype,
		runFlag,
		condition,
		command,
		outLogFile,
		errLogFile,
		label,
		profile,
		rundays,
		startWindow,
		endWindow,
		priority,
		"TEST_USER",
		"IDLE")

	if err != nil {
		return insertedTaskId, fmt.Errorf("insertTask: %w", err)
	}

	insertedTaskId, err = result.LastInsertId()
	if err != nil {
		return insertedTaskId, fmt.Errorf("insertTask: %w", err)
	}
	return insertedTaskId, nil
}

func insertStartTimes(tx *sql.Tx, id int64, startTimes []string) error {
	qry := "Insert Into sched_batch_run (task_id, start_time) Values (?,?)"

	for _, stime := range startTimes {
		if _, err := tx.Exec(qry, id, stime); err != nil {
			return fmt.Errorf("insertStartTimes: %w", er.NewDatabaseError(err.Error()))
		}
	}
	return nil
}

func insertStartMins(tx *sql.Tx, id int64, startMins []uint8) error {
	qry := "Insert Into sched_window_run (task_id, start_min) Values (?,?)"

	for _, smin := range startMins {
		if _, err := tx.Exec(qry, id, smin); err != nil {
			return fmt.Errorf("insertStartMins: %w", er.NewDatabaseError((err.Error())))
		}
	}
	return nil
}

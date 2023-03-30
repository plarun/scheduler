package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/plarun/scheduler/api/types/entity/task"
	er "github.com/plarun/scheduler/pkg/error"
)

// UpdateTask updates the task attributes by task name for an existing task
func UpdateTask(tx *sql.Tx, tsk *task.TaskEntity) error {
	// get task id and exisiting run flags of the task
	id, runFlag, err := getRunFlag(tx, tsk.Name())
	if err != nil {
		return fmt.Errorf("UpdateTask: %w", err)
	}

	// update the task attributes
	if err := updateTaskAttr(tx, tsk); err != nil {
		return fmt.Errorf("UpdateTask: %w", err)
	}

	// task relation
	if _, ok := tsk.GetFieldCondition(); ok {
		f, _ := tsk.GetField(task.FIELD_CONDITION)
		distTasks := f.(*task.Condition).DistinctTasks()
		if err := updateTaskRelation(tx, id, distTasks); err != nil {
			return fmt.Errorf("UpdateTask: %w", err)
		}
	}

	var newRunFlag task.RunType = task.RunTypeManual

	// update the task to batch run
	if f, ok := tsk.GetFieldStartTimes(); ok {
		if err := deleteStartTimes(tx, id); err != nil {
			return fmt.Errorf("UpdateTask: %w", err)
		}

		if hasStartTimes := len(f.Value()) != 0; hasStartTimes {
			if err := insertStartTimes(tx, id, f.Value()); err != nil {
				return fmt.Errorf("UpdateTask: %w", err)
			}
			newRunFlag = task.RunTypeBatch
		}
	}

	// update the start mins of task
	if f, ok := tsk.GetFieldStartMins(); ok {
		if err := deleteStartMins(tx, id); err != nil {
			return fmt.Errorf("UpdateTask: %w", err)
		}

		if hasStartMins := len(f.Value()) != 0; hasStartMins {
			if err := insertStartMins(tx, id, f.Value()); err != nil {
				return fmt.Errorf("UpdateTask: %w", err)
			}
			newRunFlag = task.RunTypeWindow
		}
	}

	// update the run window of task
	if f, ok := tsk.GetField(task.FIELD_RUN_WINDOW); ok {
		if f.Empty() {
			if err := deleteRunWindow(tx, id); err != nil {
				return fmt.Errorf("UpdateTask: %w", err)
			}
		} else {
			f, _ := tsk.GetFieldRunWindow()
			window := f.Value()
			startWindow, endWindow := window[0], window[1]
			if err := updateRunWindow(tx, id, startWindow, endWindow); err != nil {
				return fmt.Errorf("UpdateTask: %w", err)
			}
			newRunFlag = task.RunTypeWindow
		}
	}

	// update is_batch_run flag
	if newRunFlag != task.RunType(runFlag) {
		if err := updateRunFlag(tx, id, newRunFlag); err != nil {
			return fmt.Errorf("UpdateTask: %w", err)
		}
	}

	return nil
}

// updateTaskAttr updates the task attributes which are simple and straight forward
func updateTaskAttr(tx *sql.Tx, tsk *task.TaskEntity) error {
	var columns []string

	if f, ok := tsk.GetFieldMachine(); ok {
		columns = append(columns, fmt.Sprintf("machine = '%s'", f.Value()))
	}
	if f, ok := tsk.GetFieldCommand(); ok {
		columns = append(columns, fmt.Sprintf("command = '%s'", f.Value()))
	}
	if f, ok := tsk.GetFieldOutLogFile(); ok {
		columns = append(columns, fmt.Sprintf("std_out_log = '%s'", f.Value()))
	}
	if f, ok := tsk.GetFieldErrLogFile(); ok {
		columns = append(columns, fmt.Sprintf("std_err_log = '%s'", f.Value()))
	}
	if f, ok := tsk.GetFieldLabel(); ok {
		columns = append(columns, fmt.Sprintf("label = '%s'", f.Value()))
	}
	if f, ok := tsk.GetFieldProfile(); ok {
		columns = append(columns, fmt.Sprintf("profile = '%s'", f.Value()))
	}
	if f, ok := tsk.GetFieldRunDays(); ok {
		columns = append(columns, fmt.Sprintf("run_days_bit = %d", int32(f.Value())))
	}
	if f, ok := tsk.GetFieldPriority(); ok {
		columns = append(columns, fmt.Sprintf("priority = %d", f.Value()))
	}

	columnStr := strings.Join(columns, ",")
	if len(columns) != 0 {
		if _, err := tx.Exec("update sched_task set "+columnStr+" where name=?;", tsk.Name()); err != nil {
			return fmt.Errorf("updateTaskAttr: %w", er.NewDatabaseError(err.Error()))
		}
	}
	return nil
}

func updateRunWindow(tx *sql.Tx, id int64, windowStartTime, windowEndTime string) error {

	isValid := len(windowStartTime) != 0
	var nStartWindow, nEndWindow sql.NullString

	nStartWindow = sql.NullString{String: windowStartTime, Valid: isValid}
	nEndWindow = sql.NullString{String: windowEndTime, Valid: isValid}

	qry := "Update sched_task Set start_window=?, end_window=? Where id=?"

	if _, err := tx.Exec(qry, nStartWindow, nEndWindow, id); err != nil {
		return fmt.Errorf("updateRunWindow: %w", er.NewDatabaseError(err.Error()))
	}
	return nil
}

func updateRunFlag(tx *sql.Tx, id int64, flag task.RunType) error {
	qry := "Update sched_task Set run_flag=? Where id=?"

	if _, err := tx.Exec(qry, string(flag), id); err != nil {
		return fmt.Errorf("updateRunFlag: %w", er.NewDatabaseError(err.Error()))
	}
	return nil
}

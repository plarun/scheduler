package query

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/plarun/scheduler/api/types/entity/task"
)

// UpdateJob updates the job attributes by job name for an existing job
func UpdateJob(tx *sql.Tx, tsk *task.TaskEntity) error {
	// get job id and exisiting run flags of the job
	jobId, runFlag, err := getRunFlag(tx, tsk.Name())
	if err != nil {
		return fmt.Errorf("UpdateJob: %w", err)
	}

	// update the job attributes on sched_job table
	if err := updateJobAttr(tx, tsk); err != nil {
		return fmt.Errorf("UpdateJob: %w", err)
	}

	// job relation
	if _, ok := tsk.GetFieldCondition(); ok {
		f, _ := tsk.GetField(task.FIELD_CONDITION)
		distTasks := f.(*task.Condition).DistinctTasks()
		if err := updateJobRelation(tx, jobId, distTasks); err != nil {
			return fmt.Errorf("UpdateJob: %w", err)
		}
	}

	var newRunFlag task.RunType = task.RunTypeManual

	// update the job to batch run
	if f, ok := tsk.GetFieldStartTimes(); ok {
		if err := deleteStartTimes(tx, jobId); err != nil {
			return fmt.Errorf("UpdateJob: %w", err)
		}

		if hasStartTimes := len(f.Value()) != 0; hasStartTimes {
			if err := insertStartTimes(tx, jobId, f.Value()); err != nil {
				return fmt.Errorf("UpdateJob: %w", err)
			}
			newRunFlag = task.RunTypeBatch
		}
	}

	// update the start mins of job
	if f, ok := tsk.GetFieldStartMins(); ok {
		if err := deleteStartMins(tx, jobId); err != nil {
			return fmt.Errorf("UpdateJob: %w", err)
		}

		if hasStartMins := len(f.Value()) != 0; hasStartMins {
			if err := insertStartMins(tx, jobId, f.Value()); err != nil {
				return fmt.Errorf("UpdateJob: %w", err)
			}
			newRunFlag = task.RunTypeWindow
		}
	}

	// update the run window of job
	if f, ok := tsk.GetField(task.FIELD_RUN_WINDOW); ok {
		if f.Empty() {
			if err := deleteRunWindow(tx, jobId); err != nil {
				return fmt.Errorf("UpdateJob: %w", err)
			}
		} else {
			f, _ := tsk.GetFieldRunWindow()
			window := f.Value()
			startWindow, endWindow := window[0], window[1]
			if err := updateRunWindow(tx, jobId, startWindow, endWindow); err != nil {
				return fmt.Errorf("UpdateJob: %w", err)
			}
			newRunFlag = task.RunTypeWindow
		}
	}

	// update is_batch_run flag
	if newRunFlag != task.RunType(runFlag) {
		if err := updateRunFlag(tx, jobId, newRunFlag); err != nil {
			return fmt.Errorf("UpdateJob: %w", err)
		}
	}

	return nil
}

// updateJobAttr updates the job attributes which are simple and straight forward
func updateJobAttr(tx *sql.Tx, tsk *task.TaskEntity) error {
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
		columns = append(columns, fmt.Sprintf("job_profile = '%s'", f.Value()))
	}
	if f, ok := tsk.GetFieldRunDays(); ok {
		columns = append(columns, fmt.Sprintf("run_days_bit = %d", int32(f.Value())))
	}
	if f, ok := tsk.GetFieldPriority(); ok {
		columns = append(columns, fmt.Sprintf("priority = %d", f.Value()))
	}

	columnStr := strings.Join(columns, ",")
	if len(columns) != 0 {
		_, err := tx.Exec("update sched_job set "+columnStr+" where job_name=?;", tsk.Name)

		if err != nil {
			return fmt.Errorf("updateJobAttr: %v", err)
		}
	}
	return nil
}

func updateRunWindow(tx *sql.Tx, jobId int64, windowStartTime, windowEndTime string) error {

	isValid := len(windowStartTime) != 0
	var nStartWindow, nEndWindow sql.NullString

	nStartWindow = sql.NullString{String: windowStartTime, Valid: isValid}
	nEndWindow = sql.NullString{String: windowEndTime, Valid: isValid}

	qry := "Update sched_job Set start_window=?, end_window=? Where job_id=?"

	_, err := tx.Exec(qry, nStartWindow, nEndWindow, jobId)

	if err != nil {
		return fmt.Errorf("updateRunWindow: %v", err)
	}
	return nil
}

func updateRunFlag(tx *sql.Tx, jobId int64, flag task.RunType) error {
	qry := "Update sched_job Set run_flag=? Where job_id=?"

	_, err := tx.Exec(qry, string(flag), jobId)

	if err != nil {
		return fmt.Errorf("updateJobRunType: %v", err)
	}
	return nil
}

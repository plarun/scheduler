package query

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	mysql "github.com/plarun/scheduler/internal/eventserver/db"
)

// getTaskId gets task ID by task name
func getTaskId(tx *sql.Tx, name string) (int64, error) {
	var taskId int64 = 0

	qry := "Select id From sched_task Where name=?"

	row := tx.QueryRow(qry, name)

	if err := row.Scan(&taskId); err != nil {
		if err == sql.ErrNoRows {
			return taskId, fmt.Errorf("%v task not found", name)
		}
		return taskId, fmt.Errorf("getTaskId: %v", err)
	}

	return taskId, nil
}

// getTaskIdList gets list of task ID for list of tasks by task name
func getTaskIdList(tx *sql.Tx, tasks []string) ([]int64, error) {
	var ids []int64 = make([]int64, 0)

	for _, name := range tasks {
		if id, err := getTaskId(tx, name); err != nil {
			return ids, fmt.Errorf("getTaskIdList: %v", err)
		} else {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func getRunFlag(tx *sql.Tx, name string) (int64, string, error) {
	var id int64
	var runFlag string

	qry := "Select id, run_flag From sched_task Where name=?"

	row := tx.QueryRow(qry, name)

	if err := row.Scan(&id, &runFlag); err != nil {
		if err == sql.ErrNoRows {
			return id, runFlag, fmt.Errorf("getRunFlags: task not found for id %v", id)
		}
		return id, runFlag, fmt.Errorf("getRunFlags: %v", err)
	}
	return id, runFlag, nil
}

func GetTaskCommand(id int64) (string, string, string, error) {
	var command, fout, ferr string

	db := mysql.GetDatabase()

	qry := `Select command, std_out_log, std_err_log
		From sched_task
		Where id=?`

	row := db.QueryRow(qry, id)

	if err := row.Scan(&command, &fout, &ferr); err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", fmt.Errorf("GetTaskCommand: task not found for id %v", id)
		}
		return "", "", "", fmt.Errorf("GetTaskCommand: %v", err)
	}
	return command, fout, ferr, nil
}

func SetTaskStatus(id int64, state task.State) error {
	db := mysql.GetDatabase()

	qry := `Update sched_task
		Set current_status=?
		Where id=?`

	if state.IsRunning() {
		qry = `Update sched_task
		Set current_status=?, last_start_time=current_time
		Where id=?`
	} else if state.IsAborted() || state.IsFailure() || state.IsSuccess() || state.IsFrozen() {
		qry = `Update sched_task
		Set current_status=?, last_end_time=current_time
		Where id=?`
	}

	if r, err := db.Exec(qry, string(state), id); err != nil {
		return fmt.Errorf("SetTaskStatus: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetTaskStatus: %d - task id set to status %s", id, string(state))
	}
	return nil
}

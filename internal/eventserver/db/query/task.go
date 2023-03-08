package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	mysql "github.com/plarun/scheduler/internal/eventserver/db"
)

// taskExists checks whether a task is available in database
// func taskExists(name string) (bool, error) {
// 	db := db.GetDatabase()
// 	var isExists int

// 	qry := "Select Exists(Select 1 From sched_task Where name=?)"
// 	row := db.QueryRow(qry, name)

// 	err := row.Scan(&isExists)
// 	if err != nil && err != sql.ErrNoRows {
// 		return false, fmt.Errorf("taskExists: %v", err)
// 	}
// 	return isExists == 1, nil
// }

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

	if err := row.Scan(&id, &command, &fout, &ferr); err != nil {
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

	if _, err := db.Exec(qry, string(state), id); err != nil {
		return fmt.Errorf("SetTaskStatus: %v", err)
	}
	return nil
}

package query

import (
	"database/sql"
	"fmt"
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
			return id, runFlag, fmt.Errorf("task not found for id %v", id)
		}
		return id, runFlag, fmt.Errorf("getRunFlags: %v", err)
	}
	return id, runFlag, nil
}

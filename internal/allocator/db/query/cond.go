package query

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db"
)

// GetStartCondition gets the starting condition of task
func GetStartCondition(id int64) (string, error) {
	db := db.GetDatabase()

	qry := `Select start_condition 
	From sched_task 
	Where id=?`

	row := db.DB.QueryRow(qry, id)

	var cond sql.NullString
	if err := row.Scan(&cond); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("GetStartCondition: task not exist")
		}
		return "", fmt.Errorf("GetStartCondition: %w", err)
	}

	if !cond.Valid {
		return "", nil
	}
	return cond.String, nil
}

// GetPrerequisitesTaskStatus gets the distinct of tasks in the start
// condition of given task along with its current run status.
func GetDependentTasksStatus(id int64) ([]*task.TaskStatus, error) {
	db := db.GetDatabase()

	qry := `Select id, name, current_status
		From sched_task
		Where id In (
			Select cond_task_id
			From sched_task_relation
			Where task_id=?
		)`

	res := make([]*task.TaskStatus, 0)

	rows, err := db.DB.Query(qry, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("GetDependentTasksStatus: %w", err)
	}

	for rows.Next() {
		var taskId int64
		var name, status string

		rows.Scan(&taskId, &name, &status)

		ts := task.NewTaskStatus(taskId, name, task.State(status))
		res = append(res, ts)
	}

	if len(res) > 0 {
		log.Printf("GetDependentTasksStatus: %v", res)
	}
	return res, nil
}

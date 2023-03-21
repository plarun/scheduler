package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// GetStartCondition gets the starting condition of task
func GetStartCondition(id int64) (string, error) {
	db := mysql.GetDatabase()

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
	db := mysql.GetDatabase()

	qry := `Select t.id, t.name, t.current_status
		From sched_task t, sched_task_relation r
		Where t.id=r.task_id And r.cond_task_id=?`

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

		rows.Scan(&taskId)
		rows.Scan(&name)
		rows.Scan(&status)

		ts := task.NewTaskStatus(taskId, name, task.State(status))
		res = append(res, ts)
	}
	return res, nil
}

package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	mysql "github.com/plarun/scheduler/internal/eventserver/db"
)

func GetDependentTasksStatus(id int64) ([]*task.TaskStatus, error) {
	db := mysql.GetDatabase()

	qry := `Select t.id, t.name, t.current_status
		From sched_task t, sched_task_relation r
		Where t.id=r.task_id And r.cond_task_id=?`

	res := make([]*task.TaskStatus, 0)

	rows, err := db.Query(qry, id)
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

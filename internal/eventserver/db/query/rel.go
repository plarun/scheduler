package query

import (
	"database/sql"
	"fmt"

	mysql "github.com/plarun/scheduler/internal/eventserver/db"
)

func GetDependentTasks(id int64) ([]int64, error) {
	db := mysql.GetDatabase()

	qry := `Select task_id
		From sched_task_relation
		Where cond_task_id=?`

	res := make([]int64, 0)

	rows, err := db.Query(qry, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("GetDependentTasks: %w", err)
	}

	for rows.Next() {
		var taskId int64
		rows.Scan(&taskId)
		res = append(res, taskId)
	}
	return res, nil
}

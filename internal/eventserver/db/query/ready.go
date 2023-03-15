package query

import (
	"database/sql"
	"fmt"

	mysql "github.com/plarun/scheduler/internal/eventserver/db"
)

// deleteTaskRelation removes relation between task and
// tasks in starting condition. In result there will be no
// tasks in starting condition.
func PullReadyTasks() ([]int64, error) {
	db := mysql.GetDatabase()

	qry := `Select task_id
		From sched_ready
		Where lock_flag=1`

	res := make([]int64, 0)

	rows, err := db.Query(qry)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("pullReadyTasks: %w", err)
	}

	for rows.Next() {
		var taskId int64
		rows.Scan(&taskId)
		res = append(res, taskId)
	}
	return res, nil
}

func SwitchLockReadyTasks(from, to int) error {
	db := mysql.GetDatabase()

	qry := `Update sched_ready
	Set lock_flag=?
	Where lock_flag=?`

	if _, err := db.Exec(qry, to, from); err != nil {
		return fmt.Errorf("SwitchLockReadyTasks: %v", err)
	}
	return nil
}

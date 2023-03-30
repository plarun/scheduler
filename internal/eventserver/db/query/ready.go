package query

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/internal/eventserver/db"
	er "github.com/plarun/scheduler/pkg/error"
)

// deleteTaskRelation removes relation between task and
// tasks in starting condition. In result there will be no
// tasks in starting condition.
func PullReadyTasks() ([]int64, error) {
	db := db.GetDatabase()

	qry := `Select task_id
		From sched_ready
		Where lock_flag=1`

	res := make([]int64, 0)

	rows, err := db.Query(qry)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("pullReadyTasks: %w", er.NewDatabaseError(err.Error()))
	}

	for rows.Next() {
		var taskId int64
		rows.Scan(&taskId)
		res = append(res, taskId)
	}
	return res, nil
}

func SwitchLockReadyTasks(from, to int) error {
	db := db.GetDatabase()

	qry := `Update sched_ready
	Set lock_flag=?
	Where lock_flag=?`

	if _, err := db.Exec(qry, to, from); err != nil {
		return fmt.Errorf("SwitchLockReadyTasks: %w", er.NewDatabaseError(err.Error()))
	}
	return nil
}

func RemoveFromReady(id int64) error {
	db := db.GetDatabase()

	qry := `Delete From sched_ready
	Where task_id=?`

	if r, err := db.Exec(qry, id); err != nil {
		return fmt.Errorf("UnstageTask: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("RemoveFromReady: %d - task id removed from sched_ready", id)
	}
	return nil
}

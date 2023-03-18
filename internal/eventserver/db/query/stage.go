package query

import (
	"fmt"
	"log"

	mysql "github.com/plarun/scheduler/internal/eventserver/db"
)

func UnstageTask(id int64) error {
	db := mysql.GetDatabase()

	qry := `Delete From sched_stage
	Where task_id=?`

	if r, err := db.Exec(qry, id); err != nil {
		return fmt.Errorf("UnstageTask: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("UnstageTask: %d - task id removed from sched_stage", id)
	}
	return nil
}

func UnlockUnstagedTask(id int64) error {
	db := mysql.GetDatabase()

	qry := `Update sched_task
		Set lock_flag=0
		Where id=?`

	if r, err := db.Exec(qry, id); err != nil {
		return fmt.Errorf("UnlockUnstagedTask: %w", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("UnlockUnstagedTask: %d - task id is unlocked after unstaged", n)
	}
	return nil
}

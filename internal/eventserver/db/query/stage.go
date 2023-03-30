package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/eventserver/db"
	er "github.com/plarun/scheduler/pkg/error"
)

func UnstageTask(id int64) error {
	db := db.GetDatabase()

	qry := `Delete From sched_stage
	Where task_id=?`

	if r, err := db.Exec(qry, id); err != nil {
		return fmt.Errorf("UnstageTask: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("UnstageTask: %d - task id removed from sched_stage", id)
	}
	return nil
}

func UnlockUnstagedTask(id int64) error {
	db := db.GetDatabase()

	qry := `Update sched_task
		Set lock_flag=0
		Where id=?`

	if r, err := db.Exec(qry, id); err != nil {
		return fmt.Errorf("UnlockUnstagedTask: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("UnlockUnstagedTask: %d - task id is unlocked after unstaged", n)
	}
	return nil
}

func HasStagedSiblings(id int64) (bool, error) {
	db := db.GetDatabase()

	qry := `Select Count(t.id)
		From sched_task t, sched_stage s
		Where parent_id In (
			Select parent_id
			From sched_task
			Where id=?) And id<>?
		And t.id = s.task_id`

	row := db.QueryRow(qry, id, id)

	var n int
	if err := row.Scan(&n); err != nil {
		return false, fmt.Errorf("hasStagedSibling: %w", er.NewDatabaseError(err.Error()))
	}
	return n > 0, nil
}

func BundleAndSiblingsStatus(id int64) (bool, int64, int64, error) {
	db := db.GetDatabase()

	qry := `Select
		parent_id, 
		(
			Select count(id) 
			From sched_task 
			Where parent_id In (
				Select parent_id 
				From sched_task 
				Where id=?
			) And current_status In (?,?)
			And id = ?
		)
	From sched_task where id=?`

	row := db.QueryRow(qry, id, string(task.StateFailure), string(task.StateAborted), id, id)

	var parentId, n int64
	if err := row.Scan(&parentId, &n); err != nil {
		return false, 0, 0, fmt.Errorf("BundleAndSiblingsStatus: %w", er.NewDatabaseError(err.Error()))
	}

	return true, parentId, n, nil
}

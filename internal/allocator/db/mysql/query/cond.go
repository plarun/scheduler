package query

import (
	"database/sql"
	"fmt"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// GetStartCondition gets the starting condition of task
func GetStartCondition(id int) (string, error) {
	db := mysql.GetDatabase()

	qry := `Select start_condition 
	From sched_task 
	Where id=?`

	row := db.DB.QueryRow(qry, id)

	var cond string
	if err := row.Scan(&cond); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("GetStartCondition: task not exist")
		}
		return "", fmt.Errorf("GetStartCondition: %w", err)
	}

	return cond, nil
}

// GetPrerequisitesTaskStatus gets the distinct of tasks in the start
// condition of given task along with its current run status.
func GetPrerequisitesTaskStatus(id int) (map[string]string, error) {
	db := mysql.GetDatabase()

	qry := `With cond As (
		Select r.cond_task_id 
		From sched_task_relation r, sched_task t 
		Where t.id = r.task_id 
			And t.id=?
	) 
	Select t.name, t.current_status 
	From sched_task t, cond c 
	Where t.id=c.cond_task_id`

	res := make(map[string]string)

	rows, err := db.DB.Query(qry, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("GetPrerequisitesTaskStatus: %w", err)
	}

	for rows.Next() {
		var tsk, status string
		rows.Scan(&tsk)
		rows.Scan(&status)
		res[tsk] = status
	}

	return res, nil
}

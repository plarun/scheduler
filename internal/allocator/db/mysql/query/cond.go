package query

import (
	"database/sql"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func GetPrerequisitesTaskStatus(name string) (map[string]string, error) {
	db := mysql.GetDatabase()

	qry := `With cond As (
		Select r.cond_task_id 
		From sched_task_relation r, sched_task t 
		Where t.id = r.task_id 
			And t.name='?'
	) 
	Select t.name, t.current_status 
	From sched_task t, cond c 
	Where t.id=c.cond_task_id`

	res := make(map[string]string)

	rows, err := db.DB.Query(qry, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
	}

	for rows.Next() {
		var tsk, status string
		rows.Scan(&tsk)
		rows.Scan(&status)
		res[tsk] = status
	}

	return res, nil
}

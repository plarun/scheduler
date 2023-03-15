package query

import (
	"fmt"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

func setTaskstatus(id int, state task.State) error {
	db := mysql.GetDatabase()

	qry := `Update sched_task
		Set current_status=?
	Where id=?`

	_, err := db.DB.Exec(qry, string(state), id)
	if err != nil {
		return fmt.Errorf("setTaskstatus: failed to update status of task: %v", err)
	}
	return nil
}

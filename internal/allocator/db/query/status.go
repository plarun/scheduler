package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/allocator/db"
	er "github.com/plarun/scheduler/pkg/error"
)

func setTaskstatus(id int64, state task.State) error {
	db := db.GetDatabase()

	qry := `Update sched_task
		Set current_status=?
	Where id=?`

	result, err := db.DB.Exec(qry, string(state), id)
	if err != nil {
		return fmt.Errorf("setTaskstatus: failed to update status of task: %w", er.NewDatabaseError(err.Error()))
	} else if n, _ := result.RowsAffected(); n > 0 {
		log.Printf("setTaskstatus: %d - task id set to status %s", id, string(state))
	}
	return nil
}

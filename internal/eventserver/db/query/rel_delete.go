package query

import (
	"database/sql"
	"fmt"
)

// deleteTaskRelation removes relation between task and
// tasks in starting condition. In result there will be no
// tasks in starting condition.
func deleteTaskRelation(tx *sql.Tx, id int64) error {
	qry := "Delete From sched_task_relation Where task_id=?"

	if _, err := tx.Exec(qry, id); err != nil {
		return fmt.Errorf("deleteTaskRelation: %v", err)
	}

	return nil
}

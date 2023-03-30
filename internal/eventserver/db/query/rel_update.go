package query

import (
	"database/sql"
	"fmt"
)

// updateTaskRelation updates the relation of given task id
// with given new distnct condition tasks list. It will remove
// the existing relations of the given task id then will insert
// new relations with given new tasks list.
func updateTaskRelation(dbTxn *sql.Tx, id int64, condTasks []string) error {
	condIds, err := getTaskIdList(dbTxn, condTasks)
	if err != nil {
		return fmt.Errorf("updateTaskRelation: %w", err)
	}

	if err := deleteTaskRelation(dbTxn, id); err != nil {
		return fmt.Errorf("updateTaskRelation: %w", err)
	}
	if err := insertTaskRelation(dbTxn, id, condIds); err != nil {
		return fmt.Errorf("updateTaskRelation: %w", err)
	}
	return nil
}

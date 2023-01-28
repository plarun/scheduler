package query

import (
	"fmt"
	"log"

	"github.com/plarun/scheduler/internal/allocator/db/mysql"
)

// QueueTasks inserts the scheduled tasks into sched_stage for staging.
// Task cannot be staged if its already staged.
func QueueTasks() (int64, error) {
	db := mysql.GetDatabase()

	// change the status of task to '

	qry := `Insert Into sched_queue (
			task_id,
			sys_entry_date,
			priority
		)
		Select task_id
		From sched_stage
		Where is_bundle=0 And flag=1
		Union
		Select t.task_id 
		From sched_task t
			Inner Join sched_stage s On (t.parent_id=s.task_id)
		Where s.flag=1 And s.is_bundle=1`

	result, err := db.DB.Exec(qry)
	if err != nil {
		return 0, fmt.Errorf("QueueTasks: failed to push tasks into queue: %v", err)
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("QueueTasks: %v", err)
	}

	log.Printf("%d tasks pushed into queue", cnt)
	return cnt, nil
}

func DequeueTask(id int) error {
	db := mysql.GetDatabase()

	db.Lock()
	defer db.Unlock()

	qry := `Delete From sched_queue Where task_id=?`

	_, err := db.DB.Exec(qry, id)
	if err != nil {
		return fmt.Errorf("DequeueTask: failed to delete task from queue: %v", err)
	}
	return nil
}

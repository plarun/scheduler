package query

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/plarun/scheduler/api/types/entity/task"
	mysql "github.com/plarun/scheduler/internal/eventserver/db"
	"github.com/plarun/scheduler/proto"
)

// getTaskId gets task ID by task name
func getTaskId(tx *sql.Tx, name string) (int64, error) {
	var taskId int64 = 0

	qry := "Select id From sched_task Where name=?"

	row := tx.QueryRow(qry, name)

	if err := row.Scan(&taskId); err != nil {
		if err == sql.ErrNoRows {
			return taskId, fmt.Errorf("%v task not found", name)
		}
		return taskId, fmt.Errorf("getTaskId: %v", err)
	}

	return taskId, nil
}

// getTaskIdList gets list of task ID for list of tasks by task name
func getTaskIdList(tx *sql.Tx, tasks []string) ([]int64, error) {
	var ids []int64 = make([]int64, 0)

	for _, name := range tasks {
		if id, err := getTaskId(tx, name); err != nil {
			return ids, fmt.Errorf("getTaskIdList: %v", err)
		} else {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func getRunFlag(tx *sql.Tx, name string) (int64, string, error) {
	var id int64
	var runFlag string

	qry := "Select id, run_flag From sched_task Where name=?"

	row := tx.QueryRow(qry, name)

	if err := row.Scan(&id, &runFlag); err != nil {
		if err == sql.ErrNoRows {
			return id, runFlag, fmt.Errorf("getRunFlags: task not found for id %v", id)
		}
		return id, runFlag, fmt.Errorf("getRunFlags: %v", err)
	}
	return id, runFlag, nil
}

func GetTaskCommand(id int64) (string, string, string, error) {
	var command, fout, ferr string

	db := mysql.GetDatabase()

	qry := `Select command, std_out_log, std_err_log
		From sched_task
		Where id=?`

	row := db.QueryRow(qry, id)

	if err := row.Scan(&command, &fout, &ferr); err != nil {
		if err == sql.ErrNoRows {
			return "", "", "", fmt.Errorf("GetTaskCommand: task not found for id %v", id)
		}
		return "", "", "", fmt.Errorf("GetTaskCommand: %v", err)
	}
	return command, fout, ferr, nil
}

func SetTaskStatus(id int64, state task.State) error {
	db := mysql.GetDatabase()

	qry := `Update sched_task
		Set current_status=?
		Where id=?`

	if state.IsRunning() {
		qry = `Update sched_task
		Set current_status=?, last_start_time=current_time, last_end_time=null
		Where id=?`
	} else if task.IsRunnable(state) {
		qry = `Update sched_task
		Set current_status=?, last_end_time=current_time
		Where id=?`
	}

	if r, err := db.Exec(qry, string(state), id); err != nil {
		return fmt.Errorf("SetTaskStatus: %v", err)
	} else if n, _ := r.RowsAffected(); n > 0 {
		log.Printf("SetTaskStatus: %d - task id set to status %s", id, string(state))
	}

	// add entry to run history
	if task.IsRunnable(state) {
		if err := AddTaskRun(id); err != nil {
			return fmt.Errorf("SetTaskStatus: %w", err)
		}
	}
	return nil
}

func GetTaskDetails(name string) (*proto.TaskDefinition, error) {
	db := mysql.GetDatabase()

	qry := `Select
		(Select p.name From sched_task p Where p.id = t.parent_id) As parent,
		t.name, t.type, t.machine, t.command, t.start_condition, t.std_out_log, t.std_err_log, t.label, t.profile, t.run_days_bit,
		t.start_window, t.end_window, priority,
		(Select Group_concat(start_time Order By start_time Asc Separator ',') From sched_batch_run b Where b.task_id=t.id) start_times,
		(Select Group_concat(start_min Order By start_min Asc Separator ',') From sched_window_run w Where w.task_id=t.id) start_mins
	From sched_task t
	Where name=?`

	row := db.QueryRow(qry, name)

	res := &proto.TaskDefinition{}

	var (
		parent         sql.NullString
		taskName       string
		taskType       string
		machine        sql.NullString
		command        sql.NullString
		startCondition sql.NullString
		stdOutLog      sql.NullString
		stdErrLog      sql.NullString
		label          sql.NullString
		profile        sql.NullString
		runDaysBit     sql.NullString
		startWindow    sql.NullString
		endWindow      sql.NullString
		priority       sql.NullString
		startTimes     sql.NullString
		startMins      sql.NullString
	)

	err := row.Scan(
		&parent,
		&taskName,
		&taskType,
		&machine,
		&command,
		&startCondition,
		&stdOutLog,
		&stdErrLog,
		&label,
		&profile,
		&runDaysBit,
		&startWindow,
		&endWindow,
		&priority,
		&startTimes,
		&startMins)
	if err != nil {
		return res, fmt.Errorf("GetTaskDetails: %w", err)
	}

	res.Name = taskName
	res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_TYPE), Value: taskType})
	if parent.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_PARENT), Value: parent.String})
	}
	if machine.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_MACHINE), Value: machine.String})
	}
	if command.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_COMMAND), Value: command.String})
	}
	if startCondition.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_CONDITION), Value: startCondition.String})
	}
	if stdOutLog.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_OUT_LOG_FILE), Value: stdOutLog.String})
	}
	if stdErrLog.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_ERR_LOG_FILE), Value: stdErrLog.String})
	}
	if label.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_LABEL), Value: label.String})
	}
	if profile.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_PROFILE), Value: profile.String})
	}
	if runDaysBit.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_RUN_DAYS), Value: runDaysBit.String})
	}
	if startWindow.Valid && endWindow.Valid {
		win := fmt.Sprintf("%s-%s", startWindow.String, endWindow.String)
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_RUN_WINDOW), Value: win})
	}
	if priority.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_PRIORITY), Value: priority.String})
	}
	if startTimes.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_START_TIMES), Value: startTimes.String})
	}
	if startMins.Valid {
		res.Params = append(res.Params, &proto.KeyValue{Key: string(task.FIELD_START_MINS), Value: startMins.String})
	}

	if taskType == "bundle" {
		childTasks, err := getChildTasks(name)
		if err != nil {
			return res, fmt.Errorf("GetTaskDetails: %w", err)
		}

		for _, ctask := range childTasks {
			if child, err := GetTaskDetails(ctask); err != nil {
				return res, fmt.Errorf("GetTaskDetails: %w", err)
			} else {
				res.ChildrenTasks = append(res.ChildrenTasks, child)
			}
		}
	}
	return res, nil
}

func getChildTasks(name string) ([]string, error) {
	db := mysql.GetDatabase()

	qry := `Select name
		From sched_task
		Where parent_id In (
			Select id
			From sched_task
			Where name=?
		)`

	res := make([]string, 0)

	rows, err := db.Query(qry, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}
		return res, fmt.Errorf("getChildTasks: %w", err)
	}

	for rows.Next() {
		var childName string
		if err := rows.Scan(&childName); err != nil {
			return nil, fmt.Errorf("getChildTasks: %w", err)
		}
		res = append(res, childName)
	}
	return res, nil
}

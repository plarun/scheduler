package query

// getPredecessors gets all the distinct jobs in
// starting condition of job by given job_id.
// func getPredecessors(tx *sql.Tx, jobId int64) ([]string, error) {
// 	var predecessors []string

// 	qry := `Select job_name From sched_job Where job_id In (
// 		Select cond_job_id From sched_job_relation Where job_id=?)`

// 	rows, err := tx.Query(qry, jobId)

// 	if err != nil {
// 		return predecessors, fmt.Errorf("getPredecessors: %v", err)
// 	}

// 	for rows.Next() {
// 		var jobName string
// 		if err := rows.Scan(&jobName); err != nil {
// 			return predecessors, fmt.Errorf("getPredecessors: %v", err)
// 		}
// 		predecessors = append(predecessors, jobName)
// 	}
// 	return predecessors, nil
// }

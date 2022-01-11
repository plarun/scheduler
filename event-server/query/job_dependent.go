package query

import (
	"database/sql"
	"fmt"
	// pb "github.com/plarun/scheduler/event-server/data"
)

// InsertJobDependent creates job dependent relation between a job and list of jobs
func InsertJobDependent(dbTxn *sql.Tx, jobSeqId int64, dependentJobSeqIds []int64) error {
	fmt.Println(jobSeqId, dependentJobSeqIds)
	for _, dependentJobSeqId := range dependentJobSeqIds {
		_, err := dbTxn.Exec("insert into job_dependent (job_id, dependent_job_id, dependent_flag) values (?, ?, ?)", jobSeqId, dependentJobSeqId, "SU")
		if err != nil {
			return fmt.Errorf("insertJobDependent: %v", err)
		}
	}
	return nil
}

// DeleteJobDependent removes job dependent relation between a job and list of jobs
func DeleteJobDependent(dbTxn *sql.Tx, jobSeqId int64, dependentJobSeqIds []int64) error {
	fmt.Println(jobSeqId, dependentJobSeqIds)
	for _, dependentJobSeqId := range dependentJobSeqIds {
		_, err := dbTxn.Exec("delete from job_dependent where job_id=? and dependent_job_id=?", jobSeqId, dependentJobSeqId)
		if err != nil {
			return fmt.Errorf("deleteJobDependent: %v", err)
		}
	}
	return nil
}

// DeleteJobRelation
func DeleteJobRelation(dbTxn *sql.Tx, jobSeqId int64) error {
	_, err := dbTxn.Exec("delete from job_dependent where job_id=? or dependent_job_id=?", jobSeqId, jobSeqId)
	if err != nil {
		return fmt.Errorf("deleteJobRelation: %v", err)
	}
	return nil
}

// // GetJobDependents gets all the related jobs
// func GetJobDependents(dbTxn *sql.Tx, jobSeqId int64) ([]string, error) {

// 	var dependentJobs []string

// 	rows, err := dbTxn.Query("select job_name from job where job_seq_id in (select dependent_job_id from job_dependent where job_seq_id=?)", jobSeqId)
// 	if err != nil {
// 		return dependentJobs, err
// 	}

// 	for rows.Next() {
// 		var jobName string
// 		if err := rows.Scan(&jobName); err != nil {
// 			return dependentJobs, err
// 		}
// 		dependentJobs = append(dependentJobs, jobName)
// 	}

// 	return dependentJobs, nil
// }

// GetJobDependentsIdList gets list of dependent job id for requested job id
func GetJobDependentsIdList(dbTxn *sql.Tx, jobSeqId int64) ([]int64, error) {
	var dependentJobsIds []int64

	rows, err := dbTxn.Query("select dependent_job_id from job_dependent where job_id=?", jobSeqId)
	if err != nil {
		return dependentJobsIds, err
	}

	for rows.Next() {
		var dependentJobId int64
		if err := rows.Scan(&dependentJobId); err != nil {
			return dependentJobsIds, err
		}
		dependentJobsIds = append(dependentJobsIds, dependentJobId)
	}

	return dependentJobsIds, nil
}

// UpdateJobDependents updates dependent jobs for given job with latest conditions list
func UpdateJobDependents(dbTxn *sql.Tx, jobName string, conditions []string) error {
	var deleteList, insertList []int64

	jobSeqId, err := GetJobId(dbTxn, jobName)
	if err != nil {
		return err
	}

	var existingConditionJobsLookup map[int64]bool = make(map[int64]bool)
	existingConditionJobs, err := GetJobDependentsIdList(dbTxn, jobSeqId)
	if err != nil {
		return err
	}
	for _, jobId := range existingConditionJobs {
		existingConditionJobsLookup[jobId] = true
	}

	var newConditionJobsLookup map[int64]bool = make(map[int64]bool)
	newConditionJobs, err := GetJobIdList(dbTxn, conditions)
	if err != nil {
		return err
	}

	// condition jobs to be tagged
	for _, jobId := range newConditionJobs {
		newConditionJobsLookup[jobId] = true
		if _, ok := existingConditionJobsLookup[jobId]; !ok {
			insertList = append(insertList, jobId)
		}
	}

	// condition jobs to be untagged
	for _, jobId := range existingConditionJobs {
		if _, ok := newConditionJobsLookup[jobId]; !ok {
			deleteList = append(deleteList, jobId)
		}
	}

	// Actual update
	if err := DeleteJobDependent(dbTxn, jobSeqId, deleteList); err != nil {
		return err
	}
	if err := InsertJobDependent(dbTxn, jobSeqId, insertList); err != nil {
		return err
	}

	return nil
}

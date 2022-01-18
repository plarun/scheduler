package query

import (
	"database/sql"
	"fmt"
)

// InsertJobDependent creates job dependent relation between a job and list of jobs
func (database *Database) InsertJobDependent(dbTxn *sql.Tx, jobSeqId int64, dependentJobSeqIds []int64) error {
	database.lock.Lock()
	for _, dependentJobSeqId := range dependentJobSeqIds {
		_, err := dbTxn.Exec(
			`insert into job_dependent (job_id, dependent_job_id) 
			values (?, ?)`,
			jobSeqId,
			dependentJobSeqId)
		if err != nil {
			return fmt.Errorf("insertJobDependent: %v", err)
		}
	}
	database.lock.Unlock()

	return nil
}

// DeleteJobDependent removes job dependent relation between a job and list of jobs
func (database *Database) DeleteJobDependent(dbTxn *sql.Tx, jobSeqId int64, dependentJobSeqIds []int64) error {
	database.lock.Lock()
	for _, dependentJobSeqId := range dependentJobSeqIds {
		_, err := dbTxn.Exec(
			`delete from job_dependent 
			where job_id=? and dependent_job_id=?`,
			jobSeqId,
			dependentJobSeqId)
		if err != nil {
			return fmt.Errorf("deleteJobDependent: %v", err)
		}
	}
	database.lock.Unlock()

	return nil
}

// DeleteJobRelation
func (database *Database) DeleteJobRelation(dbTxn *sql.Tx, jobSeqId int64) error {
	database.lock.Lock()
	_, err := dbTxn.Exec(
		`delete from job_dependent 
		where job_id=? 
		or dependent_job_id=?`,
		jobSeqId,
		jobSeqId)
	database.lock.Unlock()

	if err != nil {
		return fmt.Errorf("deleteJobRelation: %v", err)
	}

	return nil
}

// GetJobDependents gets all the related jobs
func (database *Database) GetJobDependents(dbTxn *sql.Tx, jobSeqId int64) ([]string, error) {
	var dependentJobs []string

	database.lock.Lock()
	rows, err := dbTxn.Query(
		`select job_name from job 
		where job_seq_id in (
			select dependent_job_id from job_dependent 
			where job_seq_id=?
		)`,
		jobSeqId)
	database.lock.Unlock()

	if err != nil {
		return dependentJobs, err
	}

	for rows.Next() {
		var jobName string
		if err := rows.Scan(&jobName); err != nil {
			return dependentJobs, err
		}
		dependentJobs = append(dependentJobs, jobName)
	}

	return dependentJobs, nil
}

// GetJobDependentsIdList gets list of dependent job id for requested job id
func (database *Database) GetJobDependentsIdList(dbTxn *sql.Tx, jobSeqId int64) ([]int64, error) {
	var dependentJobsIds []int64

	database.lock.Lock()
	rows, err := dbTxn.Query(
		`select dependent_job_id from job_dependent 
		where job_id=?`,
		jobSeqId)
	database.lock.Unlock()

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
func (database *Database) UpdateJobDependents(dbTxn *sql.Tx, jobName string, conditions []string) error {
	var deleteList, insertList []int64

	jobSeqId, err := database.GetJobId(dbTxn, jobName)
	if err != nil {
		return err
	}

	var existingConditionJobsLookup map[int64]bool = make(map[int64]bool)
	existingConditionJobs, err := database.GetJobDependentsIdList(dbTxn, jobSeqId)
	if err != nil {
		return err
	}
	for _, jobId := range existingConditionJobs {
		existingConditionJobsLookup[jobId] = true
	}

	var newConditionJobsLookup map[int64]bool = make(map[int64]bool)
	newConditionJobs, err := database.GetJobIdList(dbTxn, conditions)
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
	if err := database.DeleteJobDependent(dbTxn, jobSeqId, deleteList); err != nil {
		return err
	}
	if err := database.InsertJobDependent(dbTxn, jobSeqId, insertList); err != nil {
		return err
	}

	return nil
}

// CheckConditions checks whether job's condition satisfied
func (database *Database) CheckConditions(dbTxn *sql.Tx, jobName string) (bool, error) {
	var unsatisfied int

	jobSeqId, err := database.GetJobId(dbTxn, jobName)
	if err != nil {
		return false, err
	}

	database.lock.Lock()
	row, err := database.DB.Query(
		`select count(*)
		from job
		where job_seq_id in (
			select dependent_job_id
			from job_dependent
			where job_id=?) 
		and status<>'SUCCESS'`,
		jobSeqId)
	database.lock.Unlock()

	if err != nil {
		return false, err
	}

	err = row.Scan(&unsatisfied)
	if err != nil {
		return false, err
	}

	return unsatisfied == 0, nil
}

package query

import (
	"database/sql"
	"fmt"
	"log"
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

	if database.verbose {
		log.Printf("InsertJobDependent creates job dependent relation between a job and list of jobs\n")
		log.Printf("JobSeqId: %v, DependentJobSeqIds: %v\n", jobSeqId, dependentJobSeqIds)
	}

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

	if database.verbose {
		log.Printf("DeleteJobDependent removes job dependent relation between a job and list of jobs\n")
		log.Printf("JobSeqId: %v, DependentJobSeqIds: %v\n", jobSeqId, dependentJobSeqIds)
	}

	return nil
}

// DeleteJobRelation removes all the relations of job to be deleted
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

	if database.verbose {
		log.Printf("DeleteJobRelation removes all the relations of job to be deleted\n")
		log.Printf("JobSeqId: %v\n", jobSeqId)
	}

	return nil
}

// GetPredecessors gets all the preceeding jobs
func (database *Database) GetPreceders(dbTxn *sql.Tx, jobSeqId int64) ([]string, error) {
	var preceders []string

	database.lock.Lock()
	rows, err := dbTxn.Query(
		`select job_name from job 
		where job_seq_id in (
			select dependent_job_id from job_dependent 
			where job_id=?
		)`,
		jobSeqId)
	database.lock.Unlock()

	if err != nil {
		return preceders, err
	}

	for rows.Next() {
		var jobName string
		if err := rows.Scan(&jobName); err != nil {
			return preceders, err
		}
		preceders = append(preceders, jobName)
	}

	if database.verbose {
		log.Printf("GetPredecessors gets all the preceeding jobs\n")
		log.Printf("JobSeqId: %v, Preceders: %v\n", jobSeqId, preceders)
	}

	return preceders, nil
}

// GetPredecessors gets Ids of all the preceeding jobs
func (database *Database) GetPrecedersIdList(dbTxn *sql.Tx, jobSeqId int64) ([]int64, error) {
	var precedersId []int64

	database.lock.Lock()
	rows, err := dbTxn.Query(
		`select dependent_job_id from job_dependent 
		where job_id=?`,
		jobSeqId)
	database.lock.Unlock()

	if err != nil {
		return precedersId, err
	}

	for rows.Next() {
		var dependentJobId int64
		if err := rows.Scan(&dependentJobId); err != nil {
			return precedersId, err
		}
		precedersId = append(precedersId, dependentJobId)
	}

	if database.verbose {
		log.Printf("GetPredecessors gets Ids of all the preceeding jobs\n")
		log.Printf("JobSeqId: %v, Preceders: %v\n", jobSeqId, precedersId)
	}

	return precedersId, nil
}

// UpdateJobDependents updates dependent jobs for given job with latest conditions list
func (database *Database) UpdateJobDependents(dbTxn *sql.Tx, jobName string, conditions []string) error {
	var deleteList, insertList []int64

	jobSeqId, err := database.GetJobId(dbTxn, jobName)
	if err != nil {
		return err
	}

	var existingConditionJobsLookup map[int64]bool = make(map[int64]bool)
	existingConditionJobs, err := database.GetPrecedersIdList(dbTxn, jobSeqId)
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

	if database.verbose {
		log.Printf("UpdateJobDependents updates dependent jobs for given job with latest conditions list\n")
	}

	return nil
}

// GetSuccessors gets all the successors
func (database *Database) GetSuccessors(dbTxn *sql.Tx, jobSeqId int64) ([]string, error) {
	var successors []string

	database.lock.Lock()
	rows, err := dbTxn.Query(
		`select job_name from job 
		where job_seq_id in (
			select job_id from job_dependent 
			where job_dependent_id=?
		)`,
		jobSeqId)
	database.lock.Unlock()

	if err != nil {
		return successors, err
	}

	for rows.Next() {
		var jobName string
		if err := rows.Scan(&jobName); err != nil {
			return successors, err
		}
		successors = append(successors, jobName)
	}

	if database.verbose {
		log.Printf("GetSuccessors gets all the successors\n")
		log.Printf("JobSeqId: %v, Successors: %v\n", jobSeqId, successors)
	}

	return successors, nil
}

// CheckConditions checks whether job's condition satisfied
func (database *Database) CheckConditions(dbTxn *sql.Tx, jobSeqId int) (bool, error) {
	var unsatisfied int

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

	if database.verbose {
		log.Printf("CheckConditions checks whether job's condition satisfied\n")
		log.Printf("Job: %v, Satisfied: %v\n", jobSeqId, unsatisfied == 0)
	}

	return unsatisfied == 0, nil
}

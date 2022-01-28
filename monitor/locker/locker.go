package locker

import (
	"fmt"
	"sync"
)

var jobLocker *Locker = nil

type Locker struct {
	lock       sync.Mutex
	lockStatus map[string]bool
}

func GetLocker() *Locker {
	if jobLocker == nil {
		jobLocker = &Locker{
			lock:       sync.Mutex{},
			lockStatus: make(map[string]bool),
		}
	}

	return jobLocker
}

func (locker *Locker) Put(jobName string) error {
	locker.lock.Lock()

	if _, ok := locker.lockStatus[jobName]; ok {
		return fmt.Errorf("job: %s is already available", jobName)
	}
	locker.lockStatus[jobName] = false

	locker.lock.Unlock()
	return nil
}

func (locker *Locker) Free(jobName string) error {
	locker.lock.Lock()

	if _, ok := locker.lockStatus[jobName]; !ok {
		return fmt.Errorf("job: %s is not available to lock", jobName)
	}
	delete(locker.lockStatus, jobName)

	locker.lock.Unlock()
	return nil
}

func (locker *Locker) Lock(jobName string) error {
	locker.lock.Lock()

	if _, ok := locker.lockStatus[jobName]; !ok {
		return fmt.Errorf("job: %s is not available to lock", jobName)
	}
	locker.lockStatus[jobName] = true

	locker.lock.Unlock()
	return nil
}

func (locker *Locker) Unlock(jobName string) error {
	locker.lock.Lock()

	if _, ok := locker.lockStatus[jobName]; !ok {
		return fmt.Errorf("job: %s is not available to unlock", jobName)
	}
	locker.lockStatus[jobName] = false

	locker.lock.Unlock()
	return nil
}

func (locker *Locker) Locked(jobName string) (bool, error) {
	locker.lock.Lock()

	if _, ok := locker.lockStatus[jobName]; !ok {
		return false, fmt.Errorf("job: %s is not available", jobName)
	}

	locker.lock.Unlock()
	return locker.lockStatus[jobName], nil
}

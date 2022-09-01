package wait

import (
	"log"
	"sync"

	pb "github.com/plarun/scheduler/picker/data"
)

var holder *ConcurrentHolder = nil

// ConcurrentHolder contains the jobs whose conditions are not yet satisfied
type ConcurrentHolder struct {
	lock   *sync.RWMutex
	Holder map[string]*pb.ReadyJob
}

func NewConcurrentHolder() *ConcurrentHolder {
	if holder == nil {
		holder = &ConcurrentHolder{
			lock:   &sync.RWMutex{},
			Holder: make(map[string]*pb.ReadyJob),
		}
	}
	return holder
}

func (holder *ConcurrentHolder) Hold(job *pb.ReadyJob) {
	holder.lock.Lock()
	defer holder.lock.Unlock()

	if _, ok := holder.Holder[job.GetJobName()]; !ok && !job.ConditionSatisfied {
		holder.Holder[job.GetJobName()] = job
	}
}

func (holder *ConcurrentHolder) Free(jobName string) *pb.ReadyJob {
	holder.lock.Lock()
	defer holder.lock.Unlock()

	var job *pb.ReadyJob = nil
	if _, ok := holder.Holder[jobName]; ok {
		job = holder.Holder[jobName]
	}
	delete(holder.Holder, jobName)
	return job
}

func (holder *ConcurrentHolder) Contains(jobName string) bool {
	holder.lock.RLock()
	defer holder.lock.RUnlock()

	var found = false
	if _, ok := holder.Holder[jobName]; ok {
		found = true
	}
	return found
}

func (holder *ConcurrentHolder) Print() {
	holder.lock.Lock()
	defer holder.lock.Unlock()

	for key, value := range holder.Holder {
		log.Printf("%s [%v]: {%v, %v, %v}\n", "Holder", key, value.GetJobName(), value.GetCommand(), value.GetConditionSatisfied())
	}
}

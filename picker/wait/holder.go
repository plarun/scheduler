package wait

import (
	"log"
	"sync"

	pb "github.com/plarun/scheduler/picker/data"
)

var holder *ConcurrentHolder = nil

// ConcurrentHolder contains the jobs whose conditions are not yet satisfied
type ConcurrentHolder struct {
	lock   *sync.Mutex
	Holder map[string]*pb.ReadyJob
}

func NewConcurrentHolder() *ConcurrentHolder {
	if holder == nil {
		holder = &ConcurrentHolder{
			lock:   &sync.Mutex{},
			Holder: make(map[string]*pb.ReadyJob),
		}
	}
	return holder
}

func (holder *ConcurrentHolder) Hold(job *pb.ReadyJob) {
	holder.lock.Lock()

	if _, ok := holder.Holder[job.GetJobName()]; !ok && !job.ConditionSatisfied {
		holder.Holder[job.GetJobName()] = job
	}

	holder.lock.Unlock()
}

func (holder *ConcurrentHolder) Free(jobName string) *pb.ReadyJob {
	holder.lock.Lock()

	var job *pb.ReadyJob = nil
	if _, ok := holder.Holder[jobName]; ok {
		job = holder.Holder[jobName]
	}
	delete(holder.Holder, jobName)

	holder.lock.Unlock()
	return job
}

func (holder *ConcurrentHolder) Contains(jobName string) bool {
	holder.lock.Lock()

	var found bool = false
	if _, ok := holder.Holder[jobName]; ok {
		found = true
	}

	holder.lock.Unlock()
	return found
}

func (holder *ConcurrentHolder) Print() {
	holder.lock.Lock()

	log.Println("Holder")
	for key, value := range holder.Holder {
		log.Printf("[%v]: {%v, %v, %v}\n", key, value.GetJobName(), value.GetCommand(), value.GetConditionSatisfied())
	}

	holder.lock.Unlock()
}

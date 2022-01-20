package wait

import (
	"sync"

	pb "github.com/plarun/scheduler/picker/data"
)

var holder *ConcurrentHolder = nil

// ConcurrentHolder contains the jobs whose conditions are not yet met
type ConcurrentHolder struct {
	lock   *sync.Mutex
	Holder map[string]*pb.Job
}

func NewConcurrentHolder() *ConcurrentHolder {
	if holder == nil {
		holder = &ConcurrentHolder{
			Holder: make(map[string]*pb.Job),
		}
	}
	return holder
}

func (holder *ConcurrentHolder) Hold(job *pb.Job) {
	holder.lock.Lock()

	if _, ok := holder.Holder[job.GetJobName()]; !ok && !job.ConditionSatisfied {
		holder.Holder[job.GetJobName()] = job
	}

	holder.lock.Unlock()
}

func (holder *ConcurrentHolder) Free(jobName string) *pb.Job {
	holder.lock.Lock()

	var job *pb.Job = nil
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

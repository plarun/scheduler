package queue

import (
	"fmt"
	"log"
	"sync"

	pb "github.com/plarun/scheduler/controller/data"
)

type ProcessJob struct {
	job  *pb.Job
	prev *ProcessJob
	next *ProcessJob
}

func (pque *ProcessQueue) newProcessJob(job *pb.Job) *ProcessJob {
	return &ProcessJob{
		job:  job,
		prev: nil,
		next: nil,
	}
}

func (job *ProcessJob) Job() *pb.Job {
	return job.job
}

type ProcessQueue struct {
	in   *ProcessJob
	out  *ProcessJob
	size uint32
}

var pQue *ConcurrentProcessQueue = nil

// ConcurrentProcessQueue is a concurrency wrapper for ProcessQueue
type ConcurrentProcessQueue struct {
	lock *sync.Mutex
	pQue *ProcessQueue
}

func GetProcessQueue() *ConcurrentProcessQueue {
	if pQue == nil {
		pQue = &ConcurrentProcessQueue{
			lock: &sync.Mutex{},
			pQue: &ProcessQueue{
				in:   nil,
				out:  nil,
				size: 0,
			},
		}
	}

	return pQue
}

func (que *ConcurrentProcessQueue) Push(data *pb.Job) error {
	que.lock.Lock()

	var node *ProcessJob = que.pQue.newProcessJob(data)
	if que.pQue.size == 0 {
		que.pQue.in = node
		que.pQue.out = node
	} else {
		node.next = que.pQue.in
		que.pQue.in.prev = node
		que.pQue.in = node
	}
	que.pQue.size++

	que.lock.Unlock()
	return nil
}

func (que *ConcurrentProcessQueue) Pop() (*ProcessJob, error) {
	que.lock.Lock()

	if que.pQue.size == 0 {
		return nil, fmt.Errorf("process queue is empty")
	}
	var node *ProcessJob = que.pQue.out
	que.pQue.out = que.pQue.out.prev
	if que.pQue.out != nil {
		que.pQue.out.next = nil
	}
	que.pQue.size--
	if que.pQue.size == 0 {
		que.pQue.in, que.pQue.out = nil, nil
	}

	que.lock.Unlock()
	return node, nil
}

func (que *ConcurrentProcessQueue) Size() uint32 {
	que.lock.Lock()
	len := que.pQue.size
	que.lock.Unlock()
	return len
}

func (que *ConcurrentProcessQueue) Print() {
	log.Print("[ ")
	curr := que.pQue.in
	for curr != nil {
		fmt.Printf("%v, ", curr.job.GetJobName())
		curr = curr.next
	}
	fmt.Println(" ]")
}

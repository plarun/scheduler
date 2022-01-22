package wait

import (
	"fmt"
	"log"
	"sync"

	pb "github.com/plarun/scheduler/picker/data"
)

// Double ended node
type WaitingJob struct {
	job  *pb.Job
	prev *WaitingJob
	next *WaitingJob
}

// newWaitingJob returns a new WaitingJob node
func (wque *WaitingQueue) newWaitingJob(data *pb.Job) *WaitingJob {
	return &WaitingJob{
		job:  data,
		prev: nil,
		next: nil,
	}
}

func (job *WaitingJob) Job() *pb.Job {
	return job.job
}

// WaitingQueue represents list of jobs waiting for next run
type WaitingQueue struct {
	in   *WaitingJob
	out  *WaitingJob
	size uint32
}

var wQue *ConcurrentWaitingQueue = nil

// ConcurrentWaitingQueue is a concurrency wrapper for WaitingQueue
type ConcurrentWaitingQueue struct {
	lock *sync.Mutex
	wQue *WaitingQueue
}

func NewWaitingQueue() *ConcurrentWaitingQueue {
	if wQue == nil {
		wQue = &ConcurrentWaitingQueue{
			lock: &sync.Mutex{},
			wQue: &WaitingQueue{
				in:   nil,
				out:  nil,
				size: 0,
			},
		}
	}
	return wQue
}

func (que *ConcurrentWaitingQueue) Push(data *pb.Job) error {
	que.lock.Lock()

	var node *WaitingJob = que.wQue.newWaitingJob(data)
	if que.wQue.size == 0 {
		que.wQue.in = node
		que.wQue.out = node
	} else {
		node.next = que.wQue.in
		que.wQue.in.prev = node
		que.wQue.in = node
	}
	que.wQue.size++

	que.lock.Unlock()
	return nil
}

func (que *ConcurrentWaitingQueue) Pop() (*WaitingJob, error) {
	que.lock.Lock()

	if que.wQue.size == 0 {
		return nil, fmt.Errorf("waiting queue is empty")
	}
	var node *WaitingJob = que.wQue.out
	que.wQue.out = que.wQue.out.prev
	if que.wQue.out != nil {
		que.wQue.out.next = nil
	}
	que.wQue.size--
	if que.wQue.size == 0 {
		que.wQue.in, que.wQue.out = nil, nil
	}

	que.lock.Unlock()
	return node, nil
}

// func (que *ConcurrentWaitingQueue) Remove(node *WaitingJob) *WaitingJob {
// 	if node == nil {
// 		return nil
// 	}

// 	if que.Size() == 1 {
// 		que.wQue.out, que.wQue.in = nil, nil
// 	} else if node == que.wQue.out {
// 		que.wQue.out = node.next
// 		que.wQue.out.prev = nil
// 	} else if node == que.wQue.in {
// 		que.wQue.in = node.prev
// 		que.wQue.in.next = nil
// 	} else {
// 		node.prev.next = node.next
// 		node.next.prev = node.prev
// 	}

// 	node.prev, node.next = nil, nil
// 	que.wQue.size--

// 	return node
// }

func (que *ConcurrentWaitingQueue) Size() uint32 {
	que.lock.Lock()
	len := que.wQue.size
	que.lock.Unlock()
	return len
}

func (que *ConcurrentWaitingQueue) Print() {
	log.Println()
	fmt.Print("[ ")
	curr := que.wQue.in
	for curr != nil {
		fmt.Printf("%v, ", curr.job)
		curr = curr.next
	}
	fmt.Println(" ]")
}

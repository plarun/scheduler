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

// WaitingQueue represents list of jobs waiting for next run
type WaitingQueue struct {
	in   *WaitingJob
	out  *WaitingJob
	size uint32
}

// newWaitingJob returns a new WaitingJob node
func (wque *WaitingQueue) newWaitingJob(data *pb.Job) *WaitingJob {
	return &WaitingJob{
		job:  data,
		prev: nil,
		next: nil,
	}
}

func (wque *WaitingQueue) push(data *pb.Job) error {
	var node *WaitingJob = wque.newWaitingJob(data)
	if wque.size == 0 {
		wque.in = node
		wque.out = node
	} else {
		node.next = wque.in
		wque.in.prev = node
		wque.in = node
	}
	wque.size++
	return nil
}

func (wque *WaitingQueue) pop() (*WaitingJob, error) {
	if wque.size == 0 {
		return nil, fmt.Errorf("waiting queue is empty")
	}
	var node *WaitingJob = wque.out
	wque.out = wque.out.prev
	if wque.out != nil {
		wque.out.next = nil
	}
	wque.size--
	if wque.size == 0 {
		wque.in, wque.out = nil, nil
	}
	return node, nil
}

// func (wque *WaitingQueue) empty() bool {
// 	return wque.size == 0
// }

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
	err := que.wQue.push(data)
	que.lock.Unlock()
	return err
}

func (que *ConcurrentWaitingQueue) Pop() (*WaitingJob, error) {
	que.lock.Lock()
	data, err := que.wQue.pop()
	que.lock.Unlock()
	return data, err
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

func (que *ConcurrentWaitingQueue) FeedBack(node *WaitingJob) error {
	que.lock.Lock()
	if node.next != nil || que.wQue.size > 1 {
		var err error
		if node, err = que.Pop(); err != nil {
			return err
		}
		if err = que.Push(node.job); err != nil {
			return err
		}
		node = nil
	}
	que.lock.Unlock()
	return nil
}

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

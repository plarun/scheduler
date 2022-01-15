package queue

import (
	"fmt"
	"log"
	"sync"
)

// Double ended node
type WaitingJob struct {
	job  interface{}
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
func (wque *WaitingQueue) newWaitingJob(data interface{}) *WaitingJob {
	return &WaitingJob{
		job:  data,
		prev: nil,
		next: nil,
	}
}

func (wque *WaitingQueue) push(data interface{}) error {
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

func (wque *WaitingQueue) pop() (interface{}, error) {
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

func (wque *WaitingQueue) empty() bool {
	return wque.size == 0
}

// ConcurrentWaitingQueue is a concurrency wrapper for WaitingQueue
type ConcurrentWaitingQueue struct {
	lock *sync.Mutex
	wQue *WaitingQueue
}

func NewWaitingQueue() *ConcurrentWaitingQueue {
	queue := &ConcurrentWaitingQueue{
		lock: &sync.Mutex{},
		wQue: &WaitingQueue{
			in:   nil,
			out:  nil,
			size: 0,
		},
	}
	return queue
}

func (cwque *ConcurrentWaitingQueue) Push(data interface{}) error {
	cwque.lock.Lock()
	err := cwque.wQue.push(data)
	cwque.lock.Unlock()
	return err
}

func (cwque *ConcurrentWaitingQueue) Pop() (interface{}, error) {
	cwque.lock.Lock()
	data, err := cwque.wQue.pop()
	cwque.lock.Unlock()
	return data, err
}

func (cwque *ConcurrentWaitingQueue) Size() uint32 {
	cwque.lock.Lock()
	len := cwque.wQue.size
	cwque.lock.Unlock()
	return len
}

func (cwque *ConcurrentWaitingQueue) Print() {
	log.Println()
	fmt.Print("[ ")
	curr := cwque.wQue.in
	for curr != nil {
		fmt.Printf("%v, ", curr.job)
		curr = curr.next
	}
	fmt.Println(" ]")
}

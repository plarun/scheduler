package queue

import (
	"fmt"
	"log"
	"sync"
)

type ProcessJob struct {
	job  interface{}
	prev *ProcessJob
	next *ProcessJob
}

type ProcessQueue struct {
	in   *ProcessJob
	out  *ProcessJob
	size uint32
}

func NewProcessQueue() *ProcessQueue {
	return &ProcessQueue{
		in:   nil,
		out:  nil,
		size: 0,
	}
}

func (pque *ProcessQueue) newProcessJob(job interface{}) *ProcessJob {
	return &ProcessJob{
		job:  job,
		prev: nil,
		next: nil,
	}
}

func (pque *ProcessQueue) push(data interface{}) error {
	var node *ProcessJob = pque.newProcessJob(data)
	if pque.size == 0 {
		pque.in = node
		pque.out = node
	} else {
		node.next = pque.in
		pque.in.prev = node
		pque.in = node
	}
	pque.size++
	return nil
}

func (pque *ProcessQueue) pop() (interface{}, error) {
	if pque.size == 0 {
		return nil, fmt.Errorf("waiting queue is empty")
	}
	var node *ProcessJob = pque.out
	pque.out = pque.out.prev
	if pque.out != nil {
		pque.out.next = nil
	}
	pque.size--
	if pque.size == 0 {
		pque.in, pque.out = nil, nil
	}
	return node, nil
}

func (pque *ProcessQueue) empty() bool {
	return pque.size == 0
}

// ConcurrentProcessQueue is a concurrency wrapper for ProcessQueue
type ConcurrentProcessQueue struct {
	lock *sync.Mutex
	pQue *ProcessQueue
}

func NewWaitingQueue() *ConcurrentProcessQueue {
	queue := &ConcurrentProcessQueue{
		lock: &sync.Mutex{},
		pQue: &ProcessQueue{
			in:   nil,
			out:  nil,
			size: 0,
		},
	}
	return queue
}

func (que *ConcurrentProcessQueue) Push(data interface{}) error {
	que.lock.Lock()
	err := que.pQue.push(data)
	que.lock.Unlock()
	return err
}

func (que *ConcurrentProcessQueue) Pop() (interface{}, error) {
	que.lock.Lock()
	data, err := que.pQue.pop()
	que.lock.Unlock()
	return data, err
}

func (que *ConcurrentProcessQueue) Size() uint32 {
	que.lock.Lock()
	len := que.pQue.size
	que.lock.Unlock()
	return len
}

func (que *ConcurrentProcessQueue) Print() {
	log.Println()
	fmt.Print("[ ")
	curr := que.pQue.in
	for curr != nil {
		fmt.Printf("%v, ", curr.job)
		curr = curr.next
	}
	fmt.Println(" ]")
}

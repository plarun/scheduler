package service

import (
	"log"
	"sync"
)

type WorkerPool struct {
	count int
	tasks chan Executable
}

func NewWorker(count int) *WorkerPool {
	return &WorkerPool{
		count: count,
		tasks: make(chan Executable, count),
	}
}

func (w *WorkerPool) Add(ex Executable) {
	w.tasks <- ex
}

func (w *WorkerPool) Start() {
	var wg sync.WaitGroup

	for i := 0; i < w.count; i++ {
		wg.Add(1)
		go worker(&wg, w.tasks)
	}

	wg.Wait()
}

func worker(wg *sync.WaitGroup, tasks <-chan Executable) {
	defer wg.Done()

	log.Println("worker created...")

	for task := range tasks {
		task.Execute()
	}
}

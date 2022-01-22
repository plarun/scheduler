package executor

import (
	"sync"

	pb "github.com/plarun/scheduler/controller/data"
)

const (
	ExecutorsCount = 5
)

type Executor struct {
	executing bool
	job       *pb.Job
}

func newExecutor() *Executor {
	return &Executor{
		executing: false,
		job:       nil,
	}
}

var executorPool *ExecutorPool = nil

type ExecutorPool struct {
	executors []*Executor
	lock      *sync.Mutex
}

func GetExecutorPool() *ExecutorPool {
	if executorPool == nil {
		var executors []*Executor
		for i := 0; i < ExecutorsCount; i++ {
			executors = append(executors, newExecutor())
		}
		executorPool = &ExecutorPool{
			executors: executors,
		}
	}

	return executorPool
}

func (epool *ExecutorPool) GetFreeExecutor() *Executor {
	var freeExecutor *Executor = nil
	epool.lock.Lock()

	for _, executor := range epool.executors {
		if !executor.executing {
			freeExecutor = executor
			freeExecutor.executing = true
		}
	}

	epool.lock.Unlock()
	return freeExecutor
}

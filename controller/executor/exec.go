package executor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	pb "github.com/plarun/scheduler/controller/data"
	"github.com/plarun/scheduler/controller/queue"
)

const ExecutorsCount = 4

type Executor struct {
	executing bool
	job       *pb.Job
	name      string
}

func newExecutor(id int) *Executor {
	return &Executor{
		executing: false,
		job:       nil,
		name:      "Executor" + strconv.Itoa(id),
	}
}

// execute starts the job
func (exe *Executor) execute(processJob *pb.Job) {
	defer func() {
		log.Printf("%s is freed\n", exe.name)
		exe.executing = false
		exe.job = nil
	}()

	if err := updateStatus(processJob.JobName, pb.NewStatus_CHANGE_RUNNING); err != nil {
		log.Printf("Executor.execute: %v\n", err)
	}

	failed := false

	fout, foutErr := os.OpenFile(processJob.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ferr, ferrErr := os.OpenFile(processJob.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if foutErr != nil {
		failed = true
	}
	if ferrErr != nil {
		failed = true
	}

	cmd := exec.Command(processJob.GetCommand())
	if fout != nil {
		cmd.Stdout = fout
	}
	if ferr != nil {
		cmd.Stderr = ferr
	}

	if err := cmd.Start(); err != nil {
		failed = true
	}

	if err := cmd.Wait(); err != nil || failed {
		log.Printf("%s is failed\n", processJob.JobName)
		if err := updateStatus(processJob.JobName, pb.NewStatus_CHANGE_FAILED); err != nil {
			log.Printf("Executor.execute: %v\n", err)
		}
	} else {
		log.Printf("%s is success\n", processJob.JobName)
		if err := updateStatus(processJob.JobName, pb.NewStatus_CHANGE_SUCCESS); err != nil {
			log.Printf("Executor.execute: %v\n", err)
		}
	}
}

var executorPool *ExecutorPool = nil

type ExecutorPool struct {
	executors []*Executor
	lock      *sync.Mutex
}

func GetExecutorPool() *ExecutorPool {
	if executorPool == nil {
		var executors []*Executor = make([]*Executor, 0)
		for i := 0; i < ExecutorsCount; i++ {
			executors = append(executors, newExecutor(i))
		}
		executorPool = &ExecutorPool{
			executors: executors,
			lock:      &sync.Mutex{},
		}
	}

	return executorPool
}

func (epool *ExecutorPool) getFreeExecutor() *Executor {
	var freeExecutor *Executor = nil
	epool.lock.Lock()

	for _, executor := range epool.executors {
		if !executor.executing {
			freeExecutor = executor
			freeExecutor.executing = true
			break
		}
	}

	epool.lock.Unlock()
	return freeExecutor
}

func (epool *ExecutorPool) Start() error {
	que := queue.GetProcessQueue()

	for ; true; time.Sleep(time.Millisecond * 100) {
		if que.Size() != 0 {
			executor := epool.getFreeExecutor()

			if executor != nil {
				processJob, err := que.Pop()

				if err != nil {
					log.Fatal(err)
				}
				log.Printf("%s: %s\n", executor.name, processJob.Job().GetJobName())

				go func() {
					executor.execute(processJob.Job())
				}()
			}
		}
	}

	return fmt.Errorf("unexpected failure in executor pool")
}

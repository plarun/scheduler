package executor

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	pb "github.com/plarun/scheduler/controller/data"
	"github.com/plarun/scheduler/controller/queue"
)

const ExecutorsCount = 4

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

// execute starts the job
func (exe *Executor) execute(processJob *pb.Job) {
	defer func() {
		exe.executing = false
	}()

	if err := updateStatus(processJob.JobName, pb.NewStatus_CHANGE_RUNNING); err != nil {
		log.Fatal(err)
	}

	fout, foutErr := os.OpenFile(processJob.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ferr, ferrErr := os.OpenFile(processJob.OutFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if foutErr != nil {
		log.Fatal(foutErr)
	}
	if ferrErr != nil {
		log.Fatal(ferrErr)
	}

	cmd := exec.Command(processJob.GetCommand())
	if fout != nil {
		cmd.Stdout = fout
	}
	if ferr != nil {
		cmd.Stderr = ferr
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Job: %s\nErr: %v", processJob.GetJobName(), err)
	}

	if err := cmd.Wait(); err != nil {
		log.Println("Job is failed")
		if err := updateStatus(processJob.JobName, pb.NewStatus_CHANGE_FAILED); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Job is success")
		if err := updateStatus(processJob.JobName, pb.NewStatus_CHANGE_SUCCESS); err != nil {
			log.Fatal(err)
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
		var executors []*Executor
		for i := 0; i < ExecutorsCount; i++ {
			executors = append(executors, newExecutor())
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
		}
	}

	epool.lock.Unlock()
	return freeExecutor
}

func (epool *ExecutorPool) Start() error {
	que := queue.GetProcessQueue()

	for ; true; time.Sleep(time.Millisecond * 500) {
		if que.Size() != 0 {
			executor := epool.getFreeExecutor()

			if executor != nil {
				processJob, err := que.Pop()

				if err != nil {
					panic(err)
				}
				log.Printf("Executor will execute the job: %s\n", processJob.Job())

				go func() {
					executor.execute(processJob.Job())
				}()
			}
		}
	}

	return fmt.Errorf("unexpected failure in executor pool")
}

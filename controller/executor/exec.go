package executor

import (
	"log"
	"os"
	"os/exec"
	"sync"

	pb "github.com/plarun/scheduler/controller/data"
	"github.com/plarun/scheduler/controller/queue"
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

	err := cmd.Start()
	if err != nil {
		log.Fatalf("Job: %s\nErr: %v", processJob.GetJobName(), err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal("Job is failed")
		if err := updateStatus(processJob.JobName, pb.NewStatus_CHANGE_FAILED); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Print("Job is success")
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

func (epool *ExecutorPool) Start(que *queue.ConcurrentProcessQueue) {
	for {
		if que.Size() != 0 {
			executor := epool.getFreeExecutor()
			if executor != nil {
				processJob, err := que.Pop()
				if err != nil {
					panic(err)
				}
				log.Printf("Executor will execute the job: %s\n", processJob.Job())
				executor.execute(processJob.Job())
			} else {
				continue
			}
		}
	}
}

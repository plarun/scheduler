package service

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/proto"
)

type Executable struct {
	taskId  int64
	command string
	outFile string
	errFile string
}

func NewExecutable(id int64, cmd, out, err string) Executable {
	return Executable{
		taskId:  id,
		command: cmd,
		outFile: out,
		errFile: err,
	}
}

func (e *Executable) setStatus(status proto.TaskStatus) error {
	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	req := &proto.TaskStatusRequest{
		TaskId: int64(e.taskId),
		Status: status,
	}

	conn := NewTaskStatusGrpcConnection(addr, req)
	if err := conn.Connect(); err != nil {
		return err
	}

	r, err := conn.Request()
	if err != nil {
		return err
	}

	// ignore response if success
	var ok bool
	if _, ok = r.(*proto.TaskStatusResponse); !ok {
		panic("invalid type")
	}

	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}

func (e *Executable) Execute() {
	log.Println("executing...", e.command)
	failed := false

	if err := e.setStatus(proto.TaskStatus_RUNNING); err != nil {
		log.Printf("Execute: Error - %v", err)
		return
	}

	fout, foutErr := os.OpenFile(e.outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if foutErr != nil {
		log.Printf("Execute: %v", foutErr)
	}

	ferr, ferrErr := os.OpenFile(e.errFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if ferrErr != nil {
		log.Printf("Execute: %v", ferrErr)
	}

	ex := strings.Split(e.command, " ")

	cmd := exec.Command(ex[0], ex[1:]...)
	if fout != nil {
		cmd.Stdout = fout
	}
	if ferr != nil {
		cmd.Stderr = ferr
	}

	if err := cmd.Start(); err != nil {
		failed = true
		log.Printf("Error - failed to execute the command %s : %v", e.command, err)
		return
	}

	if err := cmd.Wait(); err != nil || failed {
		if err := e.setStatus(proto.TaskStatus_FAILURE); err != nil {
			log.Printf("Execute: Error - %v", err)
		}
		log.Printf("Error - failed to wait for the command %s : %v", e.command, err)
		return
	}

	if err := e.setStatus(proto.TaskStatus_SUCCESS); err != nil {
		log.Printf("Execute: Error - %v", err)
		return
	}
	log.Println("executed", e.command)
}

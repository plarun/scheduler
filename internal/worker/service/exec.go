package service

import (
	"fmt"
	"os"
	"os/exec"

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
	failed := true

	e.setStatus(proto.TaskStatus_RUNNING)

	fout, foutErr := os.OpenFile(e.outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	ferr, ferrErr := os.OpenFile(e.outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if foutErr != nil || ferrErr != nil {
		e.setStatus(proto.TaskStatus_FAILURE)
	}

	cmd := exec.Command(e.command)
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
		e.setStatus(proto.TaskStatus_FAILURE)
	}

	e.setStatus(proto.TaskStatus_SUCCESS)
}

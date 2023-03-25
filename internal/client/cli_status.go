package client

import (
	"flag"
	"fmt"

	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/internal/client/conn"
	"github.com/plarun/scheduler/proto"
)

type statusCommand struct {
	task   string
	parsed bool
}

func newStatusCmd() Executer {
	return &statusCommand{}
}

func (sc *statusCommand) IsParsed() bool {
	return sc.parsed
}

func (sc *statusCommand) Parse(args []string) error {
	fs := flag.NewFlagSet(CMD_STATUS, flag.ContinueOnError)

	fs.StringVar(&sc.task, "j", "", "task name")

	fs.Parse(args)

	if sc.task == "" {
		return fmt.Errorf("missing task name")
	}

	sc.parsed = true
	return nil
}

func (sc *statusCommand) Exec() error {
	if !sc.IsParsed() {
		return ErrCommandNotParsed
	}

	req := &proto.TaskLatestStatusRequest{
		TaskName: sc.task,
	}

	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	conn := conn.NewStatusGrpcConnection(addr, req)

	if err := conn.Connect(); err != nil {
		return err
	}

	r, err := conn.Request()
	if err != nil {
		return err
	}

	var res *proto.TaskLatestStatusResponse
	var ok bool
	if res, ok = r.(*proto.TaskLatestStatusResponse); !ok {
		panic("invalid type")
	}

	// print task definition
	printTaskStatus(res.Status, true, "")

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func (sc *statusCommand) Usage() string {
	return USAGE_CMD_STATUS
}

func printTaskStatus(st *proto.TaskRunStatus, hasHeader bool, prefix string) {
	if hasHeader {
		fmt.Printf("\n%-65s %-17s %-17s %-10s\n", "Task Name", "Start Time", "End Time", "Status")
		fmt.Println("_________________________________________________________________ _________________ _________________ __________")
	}

	fmt.Printf("%-65s %-17s %-17s %-10s\n", prefix+st.TaskName, st.LastStartTime, st.LastEndTime, st.Status)

	for _, child := range st.Children {
		printTaskStatus(child, false, prefix+" ")
	}
}

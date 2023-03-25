package client

import (
	"flag"
	"fmt"

	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/internal/client/check"
	"github.com/plarun/scheduler/internal/client/conn"
	"github.com/plarun/scheduler/proto"
)

type runsCommand struct {
	task   string
	count  int
	date   string
	parsed bool
}

func newRunsCmd() Executer {
	return &runsCommand{}
}

func (rc *runsCommand) IsParsed() bool {
	return rc.parsed
}

func (rc *runsCommand) Parse(args []string) error {
	fs := flag.NewFlagSet(CMD_RUNS, flag.ContinueOnError)

	fs.IntVar(&rc.count, "c", 1, "number of runs")
	fs.StringVar(&rc.date, "d", "", "runs only on")
	fs.StringVar(&rc.task, "j", "", "task name")

	fs.Parse(args)

	if rc.task == "" {
		return fmt.Errorf("missing task name")
	}

	rc.parsed = true
	return nil
}

func (rc *runsCommand) Exec() error {
	if !rc.IsParsed() {
		return check.ErrCommandNotParsed
	}

	req := &proto.TaskRunsRequest{
		TaskName: rc.task,
		Count:    int32(rc.count),
		RunDate:  rc.date,
	}

	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	conn := conn.NewTaskRunsGrpcConnection(addr, req)

	if err := conn.Connect(); err != nil {
		return err
	}

	r, err := conn.Request()
	if err != nil {
		return err
	}

	var res *proto.TaskRunsResponse
	var ok bool
	if res, ok = r.(*proto.TaskRunsResponse); !ok {
		panic("invalid type")
	}

	// print task definition
	printTaskRuns(res.Runs)

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func (rc *runsCommand) Usage() string {
	return USAGE_CMD_RUNS
}

func printTaskRuns(runs []*proto.TaskRunStatus) {
	if len(runs) > 0 {
		fmt.Printf("\n%-65s %-17s %-17s %-10s\n", "Task Name", "Start Time", "End Time", "Status")
		fmt.Println("_________________________________________________________________ _________________ _________________ __________")
	}

	for _, r := range runs {
		printTaskStatus(r, false, "")
	}
}

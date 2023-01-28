package client

import (
	"flag"
	"fmt"
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

	fs.IntVar(&rc.count, "c", -1, "number of runs")
	fs.StringVar(&rc.date, "d", "01/01/1700", "runs only on")
	fs.StringVar(&rc.task, "j", "", "task name")

	fs.Parse(args)

	if rc.task == "" {
		return fmt.Errorf("missing task name")
	}

	return nil
}

func (rc *runsCommand) Exec() error {
	return nil
}

func (rc *runsCommand) Usage() string {
	return USAGE_CMD_RUNS
}

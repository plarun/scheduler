package client

import (
	"flag"
	"fmt"
)

type taskCommand struct {
	task   string
	parsed bool
}

func newTaskCmd() Executer {
	return &taskCommand{}
}

func (tc *taskCommand) IsParsed() bool {
	return tc.parsed
}

func (tc *taskCommand) Parse(args []string) error {
	fs := flag.NewFlagSet(CMD_TASK, flag.ContinueOnError)

	fs.StringVar(&tc.task, "j", "", "task name")

	fs.Parse(args)

	if tc.task == "" {
		return fmt.Errorf("missing task name")
	}

	tc.parsed = true
	return nil
}

func (tc *taskCommand) Exec() error {
	return nil
}

func (tc *taskCommand) Usage() string {
	return USAGE_CMD_TASK
}

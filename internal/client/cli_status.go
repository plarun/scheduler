package client

import (
	"flag"
	"fmt"
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
	return nil
}

func (sc *statusCommand) Usage() string {
	return USAGE_CMD_STATUS
}

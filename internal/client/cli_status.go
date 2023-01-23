package client

import (
	"flag"
	"fmt"
)

type statusCommand struct {
	job    string
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

	fs.StringVar(&sc.job, "j", "", "job name")

	fs.Parse(args)

	if sc.job == "" {
		return fmt.Errorf("missing job name")
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

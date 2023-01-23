package client

import (
	"flag"
	"fmt"
)

type jobCommand struct {
	job    string
	parsed bool
}

func newJobCmd() Executer {
	return &jobCommand{}
}

func (jc *jobCommand) IsParsed() bool {
	return jc.parsed
}

func (jc *jobCommand) Parse(args []string) error {
	fs := flag.NewFlagSet(CMD_JOB, flag.ContinueOnError)

	fs.StringVar(&jc.job, "j", "", "job name")

	fs.Parse(args)

	if jc.job == "" {
		return fmt.Errorf("missing job name")
	}

	jc.parsed = true
	return nil
}

func (jc *jobCommand) Exec() error {
	return nil
}

func (jc *jobCommand) Usage() string {
	return USAGE_CMD_JOB
}

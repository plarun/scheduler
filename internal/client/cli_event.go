package client

import (
	"flag"
	"fmt"
)

type eventCommand struct {
	event  string
	job    string
	parsed bool
}

func newEventCmd() Executer {
	return &eventCommand{}
}

func (ec *eventCommand) IsParsed() bool {
	return ec.parsed
}

func (ec *eventCommand) Parse(args []string) error {
	fs := flag.NewFlagSet(CMD_EVENT, flag.ContinueOnError)

	fs.StringVar(&ec.event, "e", "", "event name")
	fs.StringVar(&ec.job, "j", "", "job name")

	fs.Parse(args)

	if ec.event == "" {
		return fmt.Errorf("missing event name")
	}
	if ec.job == "" {
		return fmt.Errorf("missing job name")
	}

	ec.parsed = true
	return nil
}

func (ec *eventCommand) Exec() error {
	if !ec.IsParsed() {
		return ErrCommandNotParsed
	}
	return nil
}

func (ec *eventCommand) Usage() string {
	return USAGE_CMD_EVENT
}

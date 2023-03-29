package client

import (
	"flag"
	"fmt"

	"github.com/plarun/scheduler/api/types/task"
	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/internal/client/check"
	"github.com/plarun/scheduler/internal/client/conn"
	er "github.com/plarun/scheduler/internal/client/error"
	"github.com/plarun/scheduler/proto"
)

type eventCommand struct {
	event  string
	task   string
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
	fs.StringVar(&ec.task, "j", "", "task name")

	fs.Parse(args)

	if ec.event == "" {
		return fmt.Errorf("missing event name")
	}
	if ec.task == "" {
		return fmt.Errorf("missing task name")
	}

	ec.parsed = true
	return nil
}

func (ec *eventCommand) Exec() error {
	if !ec.IsParsed() {
		return check.ErrCommandNotParsed
	}

	var event task.SendEvent

	if ec.event == string(task.SendEventAbort) {
		event = task.SendEventAbort
	} else if ec.event == string(task.SendEventFreeze) {
		event = task.SendEventFreeze
	} else if ec.event == string(task.SendEventGreen) {
		event = task.SendEventRed
	} else if ec.event == string(task.SendEventReset) {
		event = task.SendEventReset
	} else if ec.event == string(task.SendEventStart) {
		event = task.SendEventStart
	} else {
		return er.ErrInvalidSendEvent
	}

	req := &proto.TaskEventRequest{
		TaskName: ec.task,
		Event:    string(event),
	}

	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	conn := conn.NewSendEventGrpcConnection(addr, req)

	if err := conn.Connect(); err != nil {
		return err
	}

	r, err := conn.Request()
	if err != nil {
		return err
	}

	var res *proto.TaskEventResponse
	var ok bool
	if res, ok = r.(*proto.TaskEventResponse); !ok {
		panic("invalid type")
	}

	if !res.Success {
		fmt.Println(res.Msg)
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func (ec *eventCommand) Usage() string {
	return USAGE_CMD_EVENT
}

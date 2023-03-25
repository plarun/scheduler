package client

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/internal/client/conn"
	"github.com/plarun/scheduler/proto"
)

// definitionCommand implements Executer
// it handles the schd_def sub program
type definitionCommand struct {
	file      *os.File
	onlyCheck bool
	parsed    bool
}

func newDefinitionCmd() Executer {
	return &definitionCommand{}
}

func (c *definitionCommand) IsParsed() bool {
	return c.parsed
}

func (c *definitionCommand) Parse(args []string) error {
	var filename string
	var cFlg bool

	log.Printf("command: \"%s\" args: \"%v\"", CMD_DEF, args)

	fs := flag.NewFlagSet(CMD_DEF, flag.ContinueOnError)

	fs.BoolVar(&cFlg, "c", false, "only check flag")
	fs.StringVar(&filename, "f", "", "file name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("file %s not found", filename)
	}

	c.file = file
	c.onlyCheck = cFlg
	c.parsed = true
	return nil
}

func (c *definitionCommand) Exec() error {
	if !c.IsParsed() {
		return ErrCommandNotParsed
	}
	defer c.file.Close()

	syntax := newDefinition(c.file)
	if err := syntax.Parse(); err != nil {
		return err
	}

	actions := syntax.Actions
	req := c.NewRequest(actions)

	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	grpcConnect := conn.NewDefinitionGrpcConnection(addr, c.onlyCheck, req)

	if err := grpcConnect.Connect(); err != nil {
		return err
	}

	r, err := grpcConnect.Request()
	if err != nil {
		return err
	}

	var res *proto.EntityActionResponse
	var ok bool
	if res, ok = r.(*proto.EntityActionResponse); !ok {
		panic("invalid type")
	}

	if res.Status.Success {
		if c.onlyCheck {
			fmt.Println("Successfully validated")
		} else {
			fmt.Println("Successfully processed")
		}
	} else {
		for _, msg := range res.Status.Errors {
			fmt.Println(msg)
		}
	}

	if err := grpcConnect.Close(); err != nil {
		return err
	}

	return nil
}

func (c *definitionCommand) Usage() string {
	return USAGE_CMD_DEF
}

func (c *definitionCommand) String() string {
	return fmt.Sprintf("command=%s file=%v only_check_flag=%v", CMD_DEF, c.file, c.onlyCheck)
}

func (c *definitionCommand) NewRequest(actions []Actioner) *proto.ParsedEntitiesRequest {
	req := &proto.ParsedEntitiesRequest{
		Tasks: make([]*proto.ParsedTaskEntity, 0),
	}

	for _, act := range actions {
		action := act.GetAction()
		// task entity action
		if task.Action(action).IsValid() {
			ent := act.(*TaskEntity)
			tsk := &proto.ParsedTaskEntity{
				Action: action,
				Target: act.GetTarget(),
				Fields: ent.Fields,
			}
			req.Tasks = append(req.Tasks, tsk)
		}
	}
	return req
}

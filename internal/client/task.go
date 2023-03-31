package client

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/config"
	"github.com/plarun/scheduler/internal/client/check"
	"github.com/plarun/scheduler/internal/client/conn"
	"github.com/plarun/scheduler/proto"
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
	if !tc.IsParsed() {
		return check.ErrCommandNotParsed
	}

	req := &proto.TaskDefinitionRequest{
		TaskName: tc.task,
	}

	addr := fmt.Sprintf(":%d", config.GetAppConfig().Service.EventServer.Port)
	conn := conn.NewTaskGrpcConnection(addr, req)

	if err := conn.Connect(); err != nil {
		return err
	}

	r, err := conn.Request()
	if err != nil {
		return err
	}

	var res *proto.TaskDefinitionResponse
	var ok bool
	if res, ok = r.(*proto.TaskDefinitionResponse); !ok {
		panic("invalid type")
	}

	if res.IsValid {
		// print task definition
		printTaskDefinition(res.Task)
	} else {
		fmt.Printf("Task %s not found\n", tc.task)
	}

	if err := conn.Close(); err != nil {
		return err
	}

	return nil
}

func (tc *taskCommand) Usage() string {
	return USAGE_CMD_TASK
}

func printTaskDefinition(def *proto.TaskDefinition) {
	fmt.Println()
	printTask(def, "")
	fmt.Println()
}

func printTask(def *proto.TaskDefinition, prefix string) {
	params := def.Params

	fmt.Printf("%sTask Name: %s\n", prefix, def.Name)

	taskType := params[string(task.FIELD_TYPE)]
	fmt.Printf("%sType: %s\n", prefix, taskType)

	if val, ok := params[string(task.FIELD_PARENT)]; ok {
		fmt.Printf("%sParent: %s\n", prefix, val)
	}

	if val, ok := params[string(task.FIELD_MACHINE)]; ok {
		fmt.Printf("%sMachine: %s\n", prefix, val)
	}

	if val, ok := params[string(task.FIELD_COMMAND)]; ok {
		fmt.Printf("%sCommand: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_CONDITION)]; ok {
		fmt.Printf("%sCondition: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_OUT_LOG_FILE)]; ok {
		fmt.Printf("%sOutLogFile: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_ERR_LOG_FILE)]; ok {
		fmt.Printf("%sErrLogFile: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_LABEL)]; ok {
		fmt.Printf("%sLabel: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_PROFILE)]; ok {
		fmt.Printf("%sProfile: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_RUN_DAYS)]; ok {
		v, _ := strconv.Atoi(val)
		days := runBitToString(v)
		fmt.Printf("%sRunDays: %s\n", prefix, days)
	}
	if val, ok := params[string(task.FIELD_START_TIMES)]; ok {
		fmt.Printf("%sStartTimes: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_RUN_WINDOW)]; ok {
		fmt.Printf("%sRunWindow: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_START_MINS)]; ok {
		fmt.Printf("%sStartMins: %s\n", prefix, val)
	}
	if val, ok := params[string(task.FIELD_PRIORITY)]; ok {
		fmt.Printf("%sPriority: %s\n", prefix, val)
	}

	// print child tasks
	if taskType == "bundle" {
		for _, child := range def.ChildrenTasks {
			fmt.Println()
			printTask(child, prefix+"  ")
		}
	}
}

// convert rundays bit to days string
func runBitToString(runDays int) string {
	var days []string
	bit := runDays

	if bit&int(task.RunDaySunday) != 0 {
		days = append(days, "su")
	}
	if bit&int(task.RunDayMonday) != 0 {
		days = append(days, "mo")
	}
	if bit&int(task.RunDayTuesday) != 0 {
		days = append(days, "tu")
	}
	if bit&int(task.RunDayWednesday) != 0 {
		days = append(days, "we")
	}
	if bit&int(task.RunDayThursday) != 0 {
		days = append(days, "th")
	}
	if bit&int(task.RunDayFriday) != 0 {
		days = append(days, "fr")
	}
	if bit&int(task.RunDaySaturday) != 0 {
		days = append(days, "sa")
	}

	return strings.Join(days, ",")
}

package client

const (
	CMD_DEF    string = "schd_def"
	CMD_EVENT  string = "schd_event"
	CMD_TASK   string = "schd_task"
	CMD_RUNS   string = "schd_runs"
	CMD_STATUS string = "schd_status"
)

const (
	// usage message for sub command schd_def
	USAGE_CMD_DEF string = `
Usage:
    schd_def [OPTION] FILE
Check and process the task actions in the file.

OPTION:
    -c, --only-check    dont process the file
FILE:
    -f, --file          input file containing task definition
                        and action for one or more tasks`
	// usage message for sub command schd_event
	USAGE_CMD_EVENT string = `
Usage:
    schd_event EVENT TASK
Send an event to task.

TASK:
    -j, --task=string   task name
EVENT:
    -e, --event=EVENT   event name should be one of following
                        start - starts the task
                        abort - stops the task
                        froze - change the status of task to FROZEN
                        reset - change the status of task to IDLE
                        chg_succ - change the status of task to SUCCESS
                        chg_fail - change the status of task to FAILURE`
	// usage message for sub command schd_task
	USAGE_CMD_TASK string = `
Usage:
    schd_task TASK
Print task definition.

TASK:
    -j, --task=string   task name`
	// usage message for sub command schd_runs
	USAGE_CMD_RUNS string = `
Usage:
    schd_runs [OPTION]... TASK
Display previous runs and status of the task.

OPTION:
    -c, --count=NUM     number of runs
    -d, --date=strings  only runs of given date
TASK:
    -j, --task=string    task name`
	// usage message for sub command schd_status
	USAGE_CMD_STATUS string = `
Usage:
    schd_status TASK_NAME
Display current run and status of the task.

TASK:
    -j, --task=string   task name`
)

// Based on the sub command, choose the corresponding
// service with provided arguments. Then validate the
// arguments.

// Executer interface represents the sub command
type Executer interface {
	// Parse the arguments of the sub command
	Parse(args []string) error
	// IsParsed checks whether the command is parsed
	IsParsed() bool
	// Execute the sub command
	Exec() error
	// Usage message of the command
	Usage() string
}

func New(subCommand string) Executer {
	switch subCommand {
	case CMD_DEF:
		return newDefinitionCmd()
	case CMD_EVENT:
		return newEventCmd()
	case CMD_TASK:
		return newTaskCmd()
	case CMD_RUNS:
		return newRunsCmd()
	case CMD_STATUS:
		return newStatusCmd()
	}
	return nil
}

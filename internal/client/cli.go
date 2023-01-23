package client

const (
	CMD_DEF    string = "schd_def"
	CMD_EVENT  string = "schd_event"
	CMD_JOB    string = "schd_job"
	CMD_RUNS   string = "schd_runs"
	CMD_STATUS string = "schd_status"
)

const (
	// usage message for sub command schd_def
	USAGE_CMD_DEF string = `
Usage:
    schd_def [OPTION] FILE
Check and process the job actions in the file.

OPTION:
    -c, --only-check    dont process the file
FILE:
    -f, --file          input file containing job definition
                        and action for one or more jobs`
	// usage message for sub command schd_event
	USAGE_CMD_EVENT string = `
Usage:
    schd_event EVENT JOB
Send an event to job.

JOB:
    -j, --job=string    job name
EVENT:
    -e, --event=EVENT   event name should be one of following
                        start - starts the job
                        abort - stops the job
                        froze - change the status of job to FROZEN
                        reset - change the status of job to IDLE
                        chg_succ - change the status of job to SUCCESS
                        chg_fail - change the status of job to FAILURE`
	// usage message for sub command schd_job
	USAGE_CMD_JOB string = `
Usage:
    schd_job JOB
Print job definition.

JOB:
    -j, --job=string   job name`
	// usage message for sub command schd_runs
	USAGE_CMD_RUNS string = `
Usage:
    schd_runs [OPTION]... JOB
Display previous runs and status of the job.

OPTION:
    -c, --count=NUM     number of runs
    -d, --date=strings  only runs of given date
JOB:
    -j, --job=string    job name`
	// usage message for sub command schd_status
	USAGE_CMD_STATUS string = `
Usage:
    schd_status JOB_NAME
Display current run and status of the job.

JOB:
    -j, --job=string   job name`
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
	case CMD_JOB:
		return newJobCmd()
	case CMD_RUNS:
		return newRunsCmd()
	case CMD_STATUS:
		return newStatusCmd()
	}
	return nil
}

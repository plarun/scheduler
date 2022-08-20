package builder

import (
	"fmt"
	"strings"

	"github.com/plarun/scheduler/client/model"
)

// InfoBuilder builds the job data from parsed Job info
type InfoBuilder struct {
	ParsedJil map[string]string
}

// BuildJil builds the jil data from parsed jil
// validates the data type of values in the jil
func (builder InfoBuilder) BuildJil() (model.JilData, error) {
	var jilData = model.JilData{}
	var err error

	// Jil Action
	if builder.ParsedJil["action"] == "insert" || builder.ParsedJil["action"] == "update" {
		jilData, err = builder.buildInsertOrUpdateJil()
	} else if builder.ParsedJil["action"] == "delete" {
		jilData, err = builder.buildDeleteJil()
	}

	if err != nil {
		return jilData, err
	}

	return jilData, nil
}

// buildInsertJil builds JIL of action type insert. Also validates the syntax for action type insert.
func (builder InfoBuilder) buildInsertOrUpdateJil() (model.JilData, error) {
	var jilData = model.JilData{}
	var attributeFlag int32 = 0

	// JIL Action
	if (builder.ParsedJil)["action"] == "insert" {
		jilData.Action = model.INSERT
	} else if (builder.ParsedJil)["action"] == "update" {
		jilData.Action = model.UPDATE
	}

	// Job name
	if jobName, ok := (builder.ParsedJil)["job_name"]; ok {
		if len(jobName) > 64 {
			return jilData, fmt.Errorf("length of the job_name should not be more than 64")
		}
		jilData.JobName = jobName
		attributeFlag = model.JobName | attributeFlag
	} else {
		return jilData, fmt.Errorf("job_name should not be empty")
	}

	// Job Command
	if command, ok := (builder.ParsedJil)["command"]; ok {
		jilData.Command = command
		attributeFlag = model.Command | attributeFlag
	} else if jilData.Action == model.INSERT {
		return jilData, fmt.Errorf("job command should not be empty")
	}

	// Job Conditions
	if conditions, ok := (builder.ParsedJil)["conditions"]; ok {
		var conditionJobs []string
		if strings.TrimSpace(conditions) == "" {
			jilData.Conditions = conditionJobs
		} else {
			conditionJobs = strings.Split(conditions, "&")
			for i := 0; i < len(conditionJobs); i++ {
				conditionJobs[i] = strings.TrimSpace(conditionJobs[i])
			}
			jilData.Conditions = conditionJobs
		}
		attributeFlag = model.Conditions | attributeFlag
	}

	// Job Standard output log path
	if stdOut, ok := (builder.ParsedJil)["std_out_log"]; ok {
		jilData.StdOutLog = stdOut
		attributeFlag = model.StdOut | attributeFlag
	}

	// Job Standard error log path
	if stdErr, ok := (builder.ParsedJil)["std_err_log"]; ok {
		jilData.StdErrLog = stdErr
		attributeFlag = model.StdErr | attributeFlag
	}

	// Machine
	if machine, ok := (builder.ParsedJil)["machine"]; ok {
		jilData.Machine = machine
		attributeFlag = model.Machine | attributeFlag
	}

	// Start time
	if startTimes, ok := (builder.ParsedJil)["start_times"]; ok {
		times := strings.Split(startTimes, ",")
		for _, time := range times {
			if err := parseTime(time); err == nil {
				jilData.StartTimes = startTimes
				attributeFlag = model.StartTimes | attributeFlag
			} else {
				return jilData, fmt.Errorf("start_time: %v is not allowed, only hh:mm format is allowed", startTimes)
			}
		}
	} else {
		// default start time
		jilData.StartTimes = "00:00:00"
		attributeFlag = model.StartTimes | attributeFlag
	}

	// Run Days
	if runDays, ok := (builder.ParsedJil)["run_days"]; ok {
		days := strings.Split(runDays, ",")
		if err := validRunDays(days); err != nil {
			return jilData, err
		}
		jilData.RunDays = runDays
		attributeFlag = model.RunDays | attributeFlag
	} else {
		jilData.RunDays = "su,mo,tu,we,th,fr,sa"
		attributeFlag = model.RunDays | attributeFlag
	}

	jilData.AttributeFlag = attributeFlag
	return jilData, nil
}

// buildDeleteJil builds JIL of action type delete. Also validates the syntax for action type delete.
func (builder InfoBuilder) buildDeleteJil() (model.JilData, error) {
	jilData := model.JilData{}
	var attributeFlag int32 = 0

	// JIL action
	jilData.Action = model.DELETE

	// Job name
	if jobName, ok := (builder.ParsedJil)["job_name"]; ok {
		jilData.JobName = jobName
		attributeFlag = model.JobName | attributeFlag
	} else {
		return jilData, fmt.Errorf("buildDeleteJil: job name is empty")
	}

	// Delete action should not take any other attribute
	nAttributes := len(builder.ParsedJil)
	if nAttributes > 2 || (builder.ParsedJil)["action"] == "" {
		return jilData, fmt.Errorf("buildDeleteJil: delete action syntax error")
	}

	jilData.AttributeFlag = attributeFlag
	return jilData, nil
}

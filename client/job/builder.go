package job

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/plarun/scheduler/client/model"
)

// JobInfoBuilder builds the job data from parsed Job info
type JobInfoBuilder struct {
	parsedJil map[string]string
}

// buildJil builds the jil data from parsed jil
// validates the data type of values in the jil
func (builder JobInfoBuilder) buildJil() (model.JilData, error) {
	var jilData = model.JilData{}
	var err error

	// Jil Action
	if builder.parsedJil["action"] == "insert" || builder.parsedJil["action"] == "update" {
		jilData, err = builder.buildInsertOrUpdateJil()
		if err != nil {
			return jilData, err
		}
	} else if builder.parsedJil["action"] == "delete" {
		jilData, err = builder.buildDeleteJil()
		if err != nil {
			return jilData, err
		}
	}

	return jilData, nil
}

// buildInsertJil builds JIL of action type insert. Also validates the syntax for action type insert.
func (builder JobInfoBuilder) buildInsertOrUpdateJil() (model.JilData, error) {
	var jilData = model.JilData{}
	var attributeFlag int32 = 0

	// JIL Action
	if (builder.parsedJil)["action"] == "insert" {
		jilData.Action = model.INSERT
	} else if (builder.parsedJil)["action"] == "update" {
		jilData.Action = model.UPDATE
	}

	// Job name
	if jobName, ok := (builder.parsedJil)["job_name"]; ok {
		if len(jobName) > 64 {
			return jilData, fmt.Errorf("length of the job_name should not be more than 64")
		}
		jilData.JobName = jobName
		attributeFlag = model.JobName | attributeFlag
	} else {
		return jilData, fmt.Errorf("job_name should not be empty")
	}

	// Job Command
	if command, ok := (builder.parsedJil)["command"]; ok {
		jilData.Command = command
		attributeFlag = model.Command | attributeFlag
	} else if jilData.Action == model.INSERT {
		return jilData, fmt.Errorf("job command should not be empty")
	}

	// Job Conditions
	if conditions, ok := (builder.parsedJil)["conditions"]; ok {
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
	if stdOut, ok := (builder.parsedJil)["std_out_log"]; ok {
		jilData.StdOutLog = stdOut
		attributeFlag = model.StdOut | attributeFlag
	}

	// Job Standard error log path
	if stdErr, ok := (builder.parsedJil)["std_err_log"]; ok {
		jilData.StdErrLog = stdErr
		attributeFlag = model.StdErr | attributeFlag
	}

	// Machine
	if machine, ok := (builder.parsedJil)["machine"]; ok {
		jilData.Machine = machine
		attributeFlag = model.Machine | attributeFlag
	}

	// Start time
	if startTimes, ok := (builder.parsedJil)["start_times"]; ok {
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
	if runDays, ok := (builder.parsedJil)["run_days"]; ok {
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
func (builder JobInfoBuilder) buildDeleteJil() (model.JilData, error) {
	jilData := model.JilData{}
	var attributeFlag int32 = 0

	// JIL action
	jilData.Action = model.DELETE

	// Job name
	if jobName, ok := (builder.parsedJil)["job_name"]; ok {
		jilData.JobName = jobName
		attributeFlag = model.JobName | attributeFlag
	} else {
		return jilData, fmt.Errorf("buildDeleteJil: job name is empty")
	}

	// Delete action should not take any other attribute
	nAttributes := len(builder.parsedJil)
	if nAttributes > 2 || (builder.parsedJil)["action"] == "" {
		return jilData, fmt.Errorf("buildDeleteJil: delete action syntax error")
	}

	jilData.AttributeFlag = attributeFlag
	return jilData, nil
}

// checkTimeFormat checks the format of time as hh:mm
func parseTime(time string) error {
	// regex to check time string in format hh:mm:ss
	timeRegex, _ := regexp.Compile("^([0-1]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$")
	if ok := timeRegex.MatchString(time); ok {
		return nil
	}
	return fmt.Errorf("start_time: %v is not allowed, only hh:mm:ss format is allowed", time)
}

// validRunDays checks format of run days
func validRunDays(days []string) error {
	var daysBitFlag int8
	// remove white spaces or tabs between week day
	// convert the week days to lower case
	for i := 0; i < len(days); i++ {
		days[i] = strings.ToLower(strings.Trim(days[i], " \t"))
	}
	for _, day := range days {
		repeatedDay := false
		switch day {
		case "su":
			if daysBitFlag&model.SU != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.SU | daysBitFlag
			}
		case "mo":
			if daysBitFlag&model.MO != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.MO | daysBitFlag
			}
		case "tu":
			if daysBitFlag&model.TU != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.TU | daysBitFlag
			}
		case "we":
			if daysBitFlag&model.WE != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.WE | daysBitFlag
			}
		case "th":
			if daysBitFlag&model.TH != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.TH | daysBitFlag
			}
		case "fr":
			if daysBitFlag&model.FR != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.FR | daysBitFlag
			}
		case "sa":
			if daysBitFlag&model.SA != 0 {
				repeatedDay = true
			} else {
				daysBitFlag = model.SA | daysBitFlag
			}
		default:
			return fmt.Errorf("invalid week: %v, days should be one of su,mo,tu,we,th,fr,sa", day)
		}
		if repeatedDay {
			return fmt.Errorf("week: %v is repated", day)
		}
	}
	return nil
}

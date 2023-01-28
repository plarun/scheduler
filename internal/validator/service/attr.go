package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/plarun/scheduler/api/types/condition"
	"github.com/plarun/scheduler/api/types/entity/task"
	"github.com/plarun/scheduler/internal/validator/errors"
	tm "github.com/plarun/scheduler/pkg/time"
)

const (
	REGEX_TASK_NAME   string = "^[0-9a-zA-Z_]+$"
	REGEX_HHMM_HHMMSS string = "^([0-1]?[0-9]|2[0-3]):[0-5][0-9](:[0-5][0-9])?$"
	REGEX_MINUTE      string = "^([0-5])([0-9])?$"
)

var (
	ActionInsertTask = string(task.ActionInsert)
	ActionUpdateTask = string(task.ActionUpdate)
	ActionDeleteTask = string(task.ActionDelete)

	TypeBundleTask   = string(task.TypeBundle)
	TypeCallableTask = string(task.TypeCallable)
)

func convertRunFlag(flag string) task.RunType {
	if flag == "batch" {
		return task.RunTypeBatch
	} else if flag == "window" {
		return task.RunTypeWindow
	} else {
		return task.RunTypeManual
	}
}

func checkFieldTaskName(name string) error {
	if len(name) == 0 {
		return errors.ErrNonEmptyValueRequired
	}
	if len(name) > 64 {
		return errors.ErrTaskMaxLength
	}

	exp, _ := regexp.Compile(REGEX_TASK_NAME)
	if ok := exp.MatchString(name); !ok {
		return errors.ErrTaskInvalidChar
	}
	return nil
}

// Action field of task should not be empty
func castFieldAction(action string) (string, error) {
	switch task.Action(action) {
	case task.ActionInsert:
		return ActionInsertTask, nil
	case task.ActionUpdate:
		return ActionUpdateTask, nil
	case task.ActionDelete:
		return ActionDeleteTask, nil
	case "":
		return "", errors.ErrNonEmptyValueRequired
	default:
		return "", errors.ErrInvalidActionAttr
	}
}

// Type field of task should not be empty
func castFieldType(tskType string) (string, error) {
	switch task.Type(tskType) {
	case task.TypeBundle:
		return TypeBundleTask, nil
	case task.TypeCallable:
		return TypeCallableTask, nil
	case "":
		return "", errors.ErrNonEmptyValueRequired
	default:
		return "", errors.ErrInvalidTypeAttr
	}
}

func castFieldStartCondition(cond string) (string, error) {
	clause, err := condition.Build(cond)
	if err != nil {
		return "", err
	}
	condStr := clause.String()
	return condStr, nil
}

func castFieldRundays(rundays string) (int32, error) {
	days := strings.Split(rundays, ",")

	var daysBitFlag, bit int32
	// remove white spaces or tabs between week day
	// convert the week days to lower case
	for i := 0; i < len(days); i++ {
		days[i] = strings.ToLower(strings.Trim(days[i], " \t"))
	}

	if len(days) == 1 {
		if days[0] == "" {
			return int32(task.RunDayEmpty), nil
		} else if days[0] == "all" {
			return int32(task.RunDayDaily), nil
		}
	}

	weekdays := map[string]int{"su": 0, "mo": 1, "tu": 2, "we": 3, "th": 4, "fr": 5, "sa": 6}

	for _, day := range days {
		if i, ok := weekdays[day]; ok {
			bit = 1 << i
			// same weekday is repeated more than once
			if daysBitFlag&bit != 0 {
				return -1, errors.ErrRepeatedRundaysAttr
			}
			// mark the weekday
			daysBitFlag |= bit
		} else {
			return -1, errors.ErrInvalidRundaysAttr
		}
	}

	return daysBitFlag, nil
}

// checkTimeFormat checks the format of time
func convertTimeToHHMM(t string) (string, error) {
	// regex to check time string in format hh:mm:ss or hh:mm
	exp, _ := regexp.Compile(REGEX_HHMM_HHMMSS)
	if ok := exp.MatchString(t); !ok {
		return "", fmt.Errorf("format other than hh:mm or hh:mm:ss")
	}

	c := strings.Count(t, ":")
	var layout string

	if c == 1 {
		layout, _ = tm.GetLayout("HHMM")
	} else {
		layout, _ = tm.GetLayout("HHMMSS")
	}

	parsed, err := time.Parse(layout, t)
	if err != nil {
		return "", fmt.Errorf("failed to parse time %v to layout hh:mm or hh:mm:ss", t)
	}
	return parsed.Format(layout), nil
}

// castFieldStartTimes validates an attribute 'start_times' of task
func castFieldStartTimes(times string) ([]string, error) {
	var startTimes []string

	times = strings.ReplaceAll(times, " ", "")
	if len(times) == 0 {
		return nil, errors.ErrNonEmptyValueRequired
	}

	for _, tm := range strings.Split(times, ",") {
		tm, err := convertTimeToHHMM(tm)
		if err != nil {
			return nil, errors.ErrInvalidStartTimesAttr
		}
		startTimes = append(startTimes, tm)
	}
	return startTimes, nil
}

// validateRunWindow validates an attribute 'run_window' of task
func castFieldRunWindow(window string) (string, string, error) {
	window = strings.ReplaceAll(window, " ", "")
	if len(window) == 0 {
		return "", "", errors.ErrNonEmptyValueRequired
	}

	windowTime := strings.Split(window, "-")
	if len(windowTime) != 2 {
		return "", "", errors.ErrInvalidRunWindowAttr
	}

	for i, time := range windowTime {
		tm, err := convertTimeToHHMM(time)
		if err != nil {
			return "", "", errors.ErrInvalidRunWindowAttr
		}
		windowTime[i] = tm
	}

	return windowTime[0], windowTime[1], nil
}

func castFieldStartMins(mins string) ([]uint8, error) {
	var startMins []uint8

	mins = strings.ReplaceAll(mins, " ", "")
	if len(mins) == 0 {
		return nil, errors.ErrNonEmptyValueRequired
	}

	exp, _ := regexp.Compile(REGEX_MINUTE)

	for _, minute := range strings.Split(mins, ",") {
		if ok := exp.MatchString(minute); !ok {
			return startMins, errors.ErrInvalidStartMinsAttr
		}
		min, _ := strconv.Atoi(minute)
		startMins = append(startMins, uint8(min))
	}
	return startMins, nil
}

func castFieldPriority(priority string) (int32, error) {
	if priority == "0" || priority == "low" || priority == "" {
		return 0, nil
	} else if priority == "1" || priority == "normal" {
		return 1, nil
	} else if priority == "2" || priority == "important" {
		return 2, nil
	} else if priority == "3" || priority == "critical" {
		return 3, nil
	} else {
		return 0, errors.ErrInvalidPriorityAttr
	}
}

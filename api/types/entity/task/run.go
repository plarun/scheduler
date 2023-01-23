package task

// RunType represents the type of run the task holds
type RunType string

const (
	// non schedulable tasks
	RunTypeManual RunType = "manual"
	// schedulable tasks with list of start times
	RunTypeBatch RunType = "batch"
	// schedulable tasks with run window range and start mins
	RunTypeWindow RunType = "window"
)

func (r RunType) IsManual() bool {
	return r == RunTypeManual
}

func (r RunType) IsBatch() bool {
	return r == RunTypeBatch
}

func (r RunType) IsWindow() bool {
	return r == RunTypeWindow
}

func (r RunType) Valid() bool {
	return r.IsManual() || r.IsBatch() || r.IsWindow()
}

func GetRunFlag(tsk *TaskEntity) RunType {
	if f, ok := tsk.GetStartTimes(); ok && len(f) != 0 {
		return RunTypeBatch
	}
	if _, ok := tsk.GetRunWindow(); ok {
		if _, ok := tsk.GetStartMins(); ok {
			return RunTypeWindow
		}
	}
	return RunTypeManual
}

// RunDay represents the days of week
type RunDay int32

const (
	RunDaySunday RunDay = 1 << iota
	RunDayMonday
	RunDayTuesday
	RunDayWednesday
	RunDayThursday
	RunDayFriday
	RunDaySaturday

	RunDayEmpty   RunDay = 0
	RunDayDaily   RunDay = RunDaySunday | RunDayMonday | RunDayTuesday | RunDayWednesday | RunDayThursday | RunDayFriday | RunDaySaturday
	RunDayWeekday RunDay = RunDayMonday | RunDayTuesday | RunDayWednesday | RunDayThursday | RunDayFriday
	RunDayWeekend RunDay = RunDaySunday | RunDaySaturday
)

func (r RunDay) IsEmpty() bool {
	return r == RunDayEmpty
}

func (r RunDay) IsDaily() bool {
	return r == RunDayDaily
}

func (r RunDay) IsWeekday() bool {
	return r == RunDayWeekday
}

func (r RunDay) IsWeekend() bool {
	return r == RunDayWeekend
}

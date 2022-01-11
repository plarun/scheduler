package model

// List of attributes for jil
var StaticAttributes = []string{
	"command",
	"conditions",
	"std_out_log",
	"std_err_log",
	"machine",
	"start_times",
	"run_days",
}

type JilAction int8

// Type of action will be performed on jil
const (
	UPDATE JilAction = iota
	INSERT
	DELETE
	NO_ACTION
)

// Attributes available in JIL data
type JilData struct {
	Action        JilAction
	JobName       string
	Command       string
	Conditions    []string
	StdOutLog     string
	StdErrLog     string
	Machine       string
	StartTimes    string
	RunDays       string
	AttributeFlag int32
}

// Represents the attribute flag
// If corresponding bit is set then that attribute has been set in JIL
const (
	JOB_NAME    int32 = 1 << iota
	COMMAND     int32 = 1 << iota
	CONDITIONS  int32 = 1 << iota
	STD_OUT     int32 = 1 << iota
	STD_ERR     int32 = 1 << iota
	MACHINE     int32 = 1 << iota
	START_TIMES int32 = 1 << iota
	RUN_DAYS    int32 = 1 << iota
)

// Represents the week days
const (
	SU int8 = 1 << iota
	MO int8 = 1 << iota
	TU int8 = 1 << iota
	WE int8 = 1 << iota
	TH int8 = 1 << iota
	FR int8 = 1 << iota
	SA int8 = 1 << iota
)

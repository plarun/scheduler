package model

import (
	pb "github.com/plarun/scheduler/client/proto"
)

// StaticAttributes List of attributes for jil
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
	//NO_ACTION
)

// JilData Attributes available in JIL data
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
	JobName    int32 = 1 << iota
	Command    int32 = 1 << iota
	Conditions int32 = 1 << iota
	StdOut     int32 = 1 << iota
	StdErr     int32 = 1 << iota
	Machine    int32 = 1 << iota
	StartTimes int32 = 1 << iota
	RunDays    int32 = 1 << iota
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

// EventTypeConv String to Job Event type
var EventTypeConv = map[string]pb.Event{
	"start":  pb.Event_START,
	"abort":  pb.Event_ABORT,
	"freeze": pb.Event_FREEZE,
	"reset":  pb.Event_RESET,
	"green":  pb.Event_GREEN,
}

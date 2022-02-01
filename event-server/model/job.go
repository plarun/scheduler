package model

import pb "github.com/plarun/scheduler/event-server/data"

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

// Job Status Number to string
var StatusTypeConv = map[string]pb.Status{
	"IDLE":    pb.Status_IDLE,
	"QUEUED":  pb.Status_QUEUED,
	"READY":   pb.Status_READY,
	"RUNNING": pb.Status_RUNNING,
	"SUCCESS": pb.Status_SUCCESS,
	"FAILED":  pb.Status_FAILED,
	"ABORTED": pb.Status_ABORTED,
	"FROZEN":  pb.Status_FROZEN,
}

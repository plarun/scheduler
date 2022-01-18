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
	"ID": pb.Status_IDLE,
	"QU": pb.Status_QUEUED,
	"RE": pb.Status_READY,
	"RU": pb.Status_RUNNING,
	"SU": pb.Status_SUCCESS,
	"FA": pb.Status_FAILED,
	"AB": pb.Status_ABORTED,
	"FZ": pb.Status_FROZEN,
}

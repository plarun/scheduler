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
	"IN": pb.Status_INACTIVE,
	"AC": pb.Status_ACTIVE,
	"ST": pb.Status_STARTED,
	"RU": pb.Status_RUNNING,
	"SU": pb.Status_SUCCESS,
	"FA": pb.Status_FAILURE,
	"OI": pb.Status_ONICE,
	"OH": pb.Status_ONHOLD,
	"TE": pb.Status_TERMINATED,
}

package model

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

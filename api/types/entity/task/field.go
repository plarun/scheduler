package task

type Field string

// Fields of task
const (
	FIELD_COMMAND      Field = "command"
	FIELD_CONDITION    Field = "condition"
	FIELD_ERR_LOG_FILE Field = "err_log_file"
	FIELD_LABEL        Field = "label"
	FIELD_MACHINE      Field = "machine"
	FIELD_OUT_LOG_FILE Field = "out_log_file"
	FIELD_PARENT       Field = "parent"
	FIELD_PRIORITY     Field = "priority"
	FIELD_PROFILE      Field = "profile"
	FIELD_TYPE         Field = "type"
	FIELD_RUN_DAYS     Field = "run_days"
	FIELD_RUN_WINDOW   Field = "run_window"
	FIELD_START_MINS   Field = "start_mins"
	FIELD_START_TIMES  Field = "start_times"
)

type ignore string

var validFields = map[Field]ignore{
	FIELD_COMMAND:      "",
	FIELD_CONDITION:    "",
	FIELD_ERR_LOG_FILE: "",
	FIELD_LABEL:        "",
	FIELD_MACHINE:      "",
	FIELD_OUT_LOG_FILE: "",
	FIELD_PARENT:       "",
	FIELD_PRIORITY:     "",
	FIELD_PROFILE:      "",
	FIELD_TYPE:         "",
	FIELD_RUN_DAYS:     "",
	FIELD_RUN_WINDOW:   "",
	FIELD_START_MINS:   "",
	FIELD_START_TIMES:  "",
}

func IsValidField(field string) bool {
	_, ok := validFields[Field(field)]
	return ok
}

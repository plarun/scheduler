package task

type OutLogFile struct {
	key   Field
	value string
	empty bool
}

func NewOutLogFile(val string, empty bool) *OutLogFile {
	return &OutLogFile{
		key:   FIELD_OUT_LOG_FILE,
		value: val,
	}
}

func (o *OutLogFile) UnimplementedTaskField() {}

func (o *OutLogFile) Key() Field {
	return o.key
}

func (o *OutLogFile) Name() string {
	return string(o.key)
}

func (o *OutLogFile) Empty() bool {
	return o.empty
}

func (o *OutLogFile) Value() string {
	return o.value
}

type ErrLogFile struct {
	key   Field
	value string
	empty bool
}

func NewErrLogFile(val string, empty bool) *ErrLogFile {
	return &ErrLogFile{
		key:   FIELD_ERR_LOG_FILE,
		value: val,
		empty: empty,
	}
}

func (e *ErrLogFile) UnimplementedTaskField() {}

func (e *ErrLogFile) Key() Field {
	return e.key
}

func (e *ErrLogFile) Name() string {
	return string(e.key)
}

func (e *ErrLogFile) Empty() bool {
	return e.empty
}

func (e *ErrLogFile) Value() string {
	return e.value
}

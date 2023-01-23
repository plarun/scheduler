package task

type RunWindow struct {
	key   Field
	value [2]string
	empty bool
}

func NewRunWindow(start, end string, empty bool) *RunWindow {
	return &RunWindow{
		key:   FIELD_RUN_WINDOW,
		value: [2]string{start, end},
		empty: empty,
	}
}

func (r *RunWindow) UnimplementedTaskField() {}

func (r *RunWindow) Key() Field {
	return r.key
}

func (r *RunWindow) Name() string {
	return string(r.key)
}

func (r *RunWindow) Empty() bool {
	return r.empty
}

func (r *RunWindow) Value() [2]string {
	return r.value
}

type StartMins struct {
	key   Field
	value []uint8
	empty bool
}

func NewStartMins(val []uint8, empty bool) *StartMins {
	return &StartMins{
		key:   FIELD_START_MINS,
		value: val,
		empty: empty,
	}
}

func (s *StartMins) UnimplementedTaskField() {}

func (s *StartMins) Key() Field {
	return s.key
}

func (s *StartMins) Name() string {
	return string(s.key)
}

func (s *StartMins) Empty() bool {
	return s.empty
}

func (s *StartMins) Value() []uint8 {
	return s.value
}

package task

type StartTimes struct {
	key   Field
	value []string
	empty bool
}

func NewStartTimes(val []string, empty bool) *StartTimes {
	return &StartTimes{
		key:   FIELD_START_TIMES,
		value: val,
		empty: empty,
	}
}

func (s *StartTimes) UnimplementedTaskField() {}

func (s *StartTimes) Key() Field {
	return s.key
}

func (s *StartTimes) Name() string {
	return string(s.key)
}

func (s *StartTimes) Empty() bool {
	return s.empty
}

func (s *StartTimes) Value() []string {
	return s.value
}

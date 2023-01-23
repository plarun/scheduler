package task

type RunDays struct {
	key   Field
	value RunDay
	empty bool
}

func NewRunDays(val int32, empty bool) *RunDays {
	return &RunDays{
		key:   FIELD_RUN_DAYS,
		value: RunDay(val),
		empty: empty,
	}
}

func (r *RunDays) UnimplementedTaskField() {}

func (r *RunDays) Key() Field {
	return r.key
}

func (r *RunDays) Name() string {
	return string(r.key)
}

func (r *RunDays) Empty() bool {
	return r.empty
}

func (r *RunDays) Value() RunDay {
	return r.value
}

package task

type Priority struct {
	key   Field
	value int32
	empty bool
}

func NewPriority(val int32, empty bool) *Priority {
	return &Priority{
		key:   FIELD_PRIORITY,
		value: val,
		empty: empty,
	}
}

func (p *Priority) UnimplementedTaskField() {}

func (p *Priority) Key() Field {
	return p.key
}

func (p *Priority) Name() string {
	return string(p.key)
}

func (p *Priority) Empty() bool {
	return p.empty
}

func (p *Priority) Value() int32 {
	return p.value
}

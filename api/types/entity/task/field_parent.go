package task

type Parent struct {
	key   Field
	value string
	empty bool
}

func NewParent(val string, empty bool) *Parent {
	return &Parent{
		key:   FIELD_PARENT,
		value: val,
		empty: empty,
	}
}

func (p *Parent) UnimplementedTaskField() {}

func (p *Parent) Key() Field {
	return p.key
}

func (p *Parent) Name() string {
	return string(p.key)
}

func (p *Parent) Empty() bool {
	return p.empty
}

func (p *Parent) Value() string {
	return p.value
}

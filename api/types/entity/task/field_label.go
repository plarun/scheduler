package task

type Label struct {
	key   Field
	value string
	empty bool
}

func NewLabel(val string, empty bool) *Label {
	return &Label{
		key:   FIELD_LABEL,
		value: val,
		empty: empty,
	}
}

func (l *Label) UnimplementedTaskField() {}

func (l *Label) Key() Field {
	return l.key
}

func (l *Label) Name() string {
	return string(l.key)
}

func (l *Label) Empty() bool {
	return l.empty
}

func (l *Label) Value() string {
	return l.value
}

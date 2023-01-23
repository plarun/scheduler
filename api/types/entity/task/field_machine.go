package task

type Machine struct {
	key   Field
	value string
	empty bool
}

func NewMachine(val string, empty bool) *Machine {
	return &Machine{
		key:   FIELD_MACHINE,
		value: val,
		empty: empty,
	}
}

func (m *Machine) UnimplementedTaskField() {}

func (m *Machine) Key() Field {
	return m.key
}

func (m *Machine) Name() string {
	return string(m.key)
}

func (m *Machine) Empty() bool {
	return m.empty
}

func (m *Machine) Value() string {
	return m.value
}

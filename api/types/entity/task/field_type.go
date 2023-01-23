package task

type TaskType struct {
	key   Field
	value Type
	empty bool
}

func NewType(t string, empty bool) *TaskType {
	return &TaskType{
		key:   FIELD_TYPE,
		value: Type(t),
		empty: empty,
	}
}

func (t *TaskType) UnimplementedTaskField() {}

func (t *TaskType) Key() Field {
	return t.key
}

func (t *TaskType) Name() string {
	return string(t.key)
}

func (t *TaskType) Empty() bool {
	return t.empty
}

func (t *TaskType) Value() Type {
	return t.value
}

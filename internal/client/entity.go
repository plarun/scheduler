package client

// Actioner interface type represents the entity
type Actioner interface {
	GetTarget() string
	GetAction() string
}

// TaskEntity represents the action of parsed task
type TaskEntity struct {
	Action string
	Target string
	Fields map[string]string
}

func newTaskEntity(action, target string) *TaskEntity {
	return &TaskEntity{
		Action: action,
		Target: target,
		Fields: make(map[string]string),
	}
}

func (t *TaskEntity) GetAction() string {
	return t.Action
}

func (t *TaskEntity) GetTarget() string {
	return t.Target
}

func (t *TaskEntity) AddField(key, value string) {

	t.Fields[key] = value
}

func (t *TaskEntity) FieldExist(key string) bool {
	_, ok := t.Fields[key]
	return ok
}

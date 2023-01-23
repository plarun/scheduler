package task

// Type represents type of the task
type Type string

// task types
const (
	TypeCallable Type = "callable"
	TypeBundle   Type = "bundle"
)

func (t Type) IsCallable() bool {
	return t == TypeCallable
}

func (t Type) IsBundle() bool {
	return t == TypeBundle
}

func (t Type) Valid() bool {
	return t.IsCallable() || t.IsBundle()
}

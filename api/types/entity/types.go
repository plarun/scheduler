package entity

type Type string

const (
	TypeUnknown Type = ""
	TypeTask    Type = "task"
	TypeMachine Type = "machine"
)

func (t Type) IsUnknown() bool {
	return t == TypeUnknown
}

func (t Type) IsTask() bool {
	return t == TypeTask
}

func (t Type) IsMachine() bool {
	return t == TypeMachine
}

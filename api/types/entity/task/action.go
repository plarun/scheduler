package task

// Action represents the type of action to be performed on task
type Action string

const (
	// Insert a new task with fields
	ActionInsert Action = "insert_task"
	// Update the editable fields of existing task
	ActionUpdate Action = "update_task"
	// Delete an existing task
	ActionDelete Action = "delete_task"
)

func (a Action) IsInsert() bool {
	return a == ActionInsert
}

func (a Action) IsUpdate() bool {
	return a == ActionUpdate
}

func (a Action) IsDelete() bool {
	return a == ActionDelete
}

func (a Action) IsValid() bool {
	return a.IsInsert() || a.IsUpdate() || a.IsDelete()
}

func IsValidAction(action string) bool {
	return action == string(ActionDelete) || action == string(ActionInsert) || action == string(ActionUpdate)
}

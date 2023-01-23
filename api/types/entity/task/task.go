package task

type TaskField interface {
	// To mention that a particular is part of an entity
	UnimplementedTaskField()
	// Field
	Key() Field
	// Field name
	Name() string
	// To check if field is empty
	Empty() bool
}

type TaskEntity struct {
	name   string
	fields map[Field]TaskField
}

func NewTaskEntity(name string) *TaskEntity {
	return &TaskEntity{
		name:   name,
		fields: make(map[Field]TaskField),
	}
}

func (t *TaskEntity) Name() string {
	return t.name
}

func (t *TaskEntity) HasField(field Field) bool {
	_, ok := t.fields[field]
	return ok
}

func (t *TaskEntity) GetField(field Field) (TaskField, bool) {
	if f, ok := t.fields[field]; ok {
		return f, ok
	}
	return nil, false
}

func (t *TaskEntity) SetFieldCommand(val string, empty bool) {
	t.fields[FIELD_COMMAND] = NewCommand(val, empty)
}

func (t *TaskEntity) GetFieldCommand() (*Command, bool) {
	f, ok := t.GetField(FIELD_COMMAND)
	if ok {
		return f.(*Command), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetCommand() (string, bool) {
	f, ok := t.GetFieldCommand()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldCondition(val string, empty bool) {
	t.fields[FIELD_CONDITION] = NewCondition(val, empty)
}

func (t *TaskEntity) GetFieldCondition() (*Condition, bool) {
	f, ok := t.GetField(FIELD_CONDITION)
	if ok {
		return f.(*Condition), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetCondition() (string, bool) {
	f, ok := t.GetFieldCondition()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldLabel(val string, empty bool) {
	t.fields[FIELD_LABEL] = NewLabel(val, empty)
}

func (t *TaskEntity) GetFieldLabel() (*Label, bool) {
	f, ok := t.GetField(FIELD_LABEL)
	if ok {
		return f.(*Label), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetLabel() (string, bool) {
	f, ok := t.GetFieldLabel()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldOutLogFile(val string, empty bool) {
	t.fields[FIELD_OUT_LOG_FILE] = NewOutLogFile(val, empty)
}

func (t *TaskEntity) GetFieldOutLogFile() (*OutLogFile, bool) {
	f, ok := t.GetField(FIELD_OUT_LOG_FILE)
	if ok {
		return f.(*OutLogFile), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetOutLogFile() (string, bool) {
	f, ok := t.GetFieldOutLogFile()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldErrLogFile(val string, empty bool) {
	t.fields[FIELD_ERR_LOG_FILE] = NewErrLogFile(val, empty)
}

func (t *TaskEntity) GetFieldErrLogFile() (*ErrLogFile, bool) {
	f, ok := t.GetField(FIELD_ERR_LOG_FILE)
	if ok {
		return f.(*ErrLogFile), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetErrLogFile() (string, bool) {
	f, ok := t.GetFieldErrLogFile()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldMachine(val string, empty bool) {
	t.fields[FIELD_MACHINE] = NewMachine(val, empty)
}

func (t *TaskEntity) GetFieldMachine() (*Machine, bool) {
	f, ok := t.GetField(FIELD_MACHINE)
	if ok {
		return f.(*Machine), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetMachine() (string, bool) {
	f, ok := t.GetFieldMachine()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldParent(val string, empty bool) {
	t.fields[FIELD_PARENT] = NewParent(val, empty)
}

func (t *TaskEntity) GetFieldParent() (*Parent, bool) {
	f, ok := t.GetField(FIELD_PARENT)
	if ok {
		return f.(*Parent), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetParent() (string, bool) {
	f, ok := t.GetFieldParent()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldPriority(val int32, empty bool) {
	t.fields[FIELD_PRIORITY] = NewPriority(val, empty)
}

func (t *TaskEntity) GetFieldPriority() (*Priority, bool) {
	f, ok := t.GetField(FIELD_PRIORITY)
	if ok {
		return f.(*Priority), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetPriority() (int32, bool) {
	f, ok := t.GetFieldPriority()
	if ok {
		return f.Value(), ok
	}
	return 0, ok
}

func (t *TaskEntity) SetFieldProfile(val string, empty bool) {
	t.fields[FIELD_PROFILE] = NewProfile(val, empty)
}

func (t *TaskEntity) GetFieldProfile() (*Profile, bool) {
	f, ok := t.GetField(FIELD_PROFILE)
	if ok {
		return f.(*Profile), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetProfile() (string, bool) {
	f, ok := t.GetFieldProfile()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

func (t *TaskEntity) SetFieldRunDays(val int32, empty bool) {
	t.fields[FIELD_RUN_DAYS] = NewRunDays(val, empty)
}

func (t *TaskEntity) GetFieldRunDays() (*RunDays, bool) {
	f, ok := t.GetField(FIELD_PROFILE)
	if ok {
		return f.(*RunDays), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetRunDays() (RunDay, bool) {
	f, ok := t.GetFieldRunDays()
	if ok {
		return f.Value(), ok
	}
	return RunDayEmpty, ok
}

func (t *TaskEntity) SetFieldStartTimes(val []string, empty bool) {
	t.fields[FIELD_START_TIMES] = NewStartTimes(val, empty)
}

func (t *TaskEntity) GetFieldStartTimes() (*StartTimes, bool) {
	f, ok := t.GetField(FIELD_START_TIMES)
	if ok {
		return f.(*StartTimes), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetStartTimes() ([]string, bool) {
	f, ok := t.GetFieldStartTimes()
	if ok {
		return f.Value(), ok
	}
	return nil, ok
}

func (t *TaskEntity) SetFieldRunWindow(start, end string, empty bool) {
	t.fields[FIELD_RUN_WINDOW] = NewRunWindow(start, end, empty)
}

func (t *TaskEntity) GetFieldRunWindow() (*RunWindow, bool) {
	f, ok := t.GetField(FIELD_RUN_WINDOW)
	if ok {
		return f.(*RunWindow), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetRunWindow() ([2]string, bool) {
	f, ok := t.GetFieldRunWindow()
	if ok {
		return f.Value(), ok
	}
	return [2]string{"", ""}, ok
}

func (t *TaskEntity) SetFieldStartMins(val []uint8, empty bool) {
	t.fields[FIELD_START_MINS] = NewStartMins(val, empty)
}

func (t *TaskEntity) GetFieldStartMins() (*StartMins, bool) {
	f, ok := t.GetField(FIELD_START_MINS)
	if ok {
		return f.(*StartMins), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetStartMins() ([]uint8, bool) {
	f, ok := t.GetFieldStartMins()
	if ok {
		return f.Value(), ok
	}
	return nil, ok
}

func (t *TaskEntity) SetFieldType(val string, empty bool) {
	t.fields[FIELD_TYPE] = NewType(val, empty)
}

func (t *TaskEntity) GetFieldType() (*TaskType, bool) {
	f, ok := t.GetField(FIELD_TYPE)
	if ok {
		return f.(*TaskType), ok
	}
	return nil, ok
}

func (t *TaskEntity) GetType() (Type, bool) {
	f, ok := t.GetFieldType()
	if ok {
		return f.Value(), ok
	}
	return "", ok
}

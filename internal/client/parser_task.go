package client

import (
	"strings"

	"github.com/plarun/scheduler/api/types/entity/task"
)

// TaskEntityParser represents a parser for task entity
type TaskEntityParser struct {
	Task   *TaskEntity
	Parsed bool
}

func newTaskEntityParser(action, target string) *TaskEntityParser {
	return &TaskEntityParser{
		Task: newTaskEntity(action, target),
	}
}

func (t *TaskEntityParser) IsParsed() bool {
	return t.Parsed
}

func (t *TaskEntityParser) Get() Actioner {
	return t.Task
}

func (t *TaskEntityParser) Parse(line string) error {
	// empty line
	if t.IsEmptyLine(line) {
		if !t.Parsed {
			t.Parsed = true
			// task entity should be completed
		}
		return nil
	}

	if t.IsComment(line) {
		if !t.Parsed {
			return ErrDefUnwantedComment
		}
		return nil
	}

	key, value, err := splitToKeyValue(line)
	if err != nil {
		return err
	}

	value = strings.TrimLeft(value, " \t")

	if !task.IsValidField(key) {
		return ErrDefInvalidField
	}
	t.Task.AddField(key, value)

	return nil
}

func (t *TaskEntityParser) IsComment(line string) bool {
	return line[:1] == "#"
}

func (t *TaskEntityParser) IsEmptyLine(line string) bool {
	return len(line) == 0
}

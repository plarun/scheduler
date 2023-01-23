package client

import "github.com/plarun/scheduler/api/types/entity"

// ActionParser interface type represents a parser for action definition
type ActionParser interface {
	IsParsed() bool
	Get() Actioner
	Parse(line string) error
}

func newParser(ent entity.Type, action, target string) ActionParser {
	switch ent {
	case entity.TypeTask:
		return newTaskEntityParser(action, target)
	}
	return nil
}

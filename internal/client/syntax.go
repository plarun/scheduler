package client

import (
	"bufio"
	"os"
	"strings"

	"github.com/plarun/scheduler/api/types/entity"
	"github.com/plarun/scheduler/api/types/entity/task"
)

// Definition represents a type to parse action definition of entities in the file.
// Entities can be Task, Machine, etc
type Definition struct {
	File      *os.File
	Actions   []Actioner
	LineNo    uint
	Line      string
	Parser    ActionParser
	HasParser bool
}

func newDefinition(file *os.File) *Definition {
	return &Definition{
		File:      file,
		Actions:   make([]Actioner, 0),
		Line:      "",
		LineNo:    0,
		HasParser: false,
	}
}

func (d *Definition) SetParser(parser ActionParser) {
	d.Parser = parser
	d.HasParser = true
}

func (d *Definition) Parse() error {

	scanner := bufio.NewScanner(d.File)
	for scanner.Scan() {
		d.Line = scanner.Text()
		d.LineNo++

		// remove leading and trailing spaces
		d.Line = strings.TrimSpace(d.Line)

		// set parser
		if !d.HasParser {
			if d.IsComment(d.Line) || d.IsEmptyLine(d.Line) {
				continue
			}

			// identify the action and its target
			action, target, err := splitToKeyValue(d.Line)
			if err != nil {
				return ErrDefInvalidKeyValue
			}

			ent := getEntityType(action)
			if ent.IsUnknown() {
				return ErrDefInvalidAction
			}

			parser := newParser(ent, action, target)
			d.SetParser(parser)
			continue
		}

		// Parse the line
		if err := d.Parser.Parse(d.Line); err != nil {
			return err
		}

		// Add to actions if an entity is parsed
		if d.Parser.IsParsed() {
			d.Actions = append(d.Actions, d.Parser.Get())
			// reset parser
			d.Parser = nil
			d.HasParser = false
		}
	}

	if !d.Parser.IsParsed() {
		d.Actions = append(d.Actions, d.Parser.Get())
		d.Parser = nil
		d.HasParser = false
	}

	return nil
}

func (d *Definition) IsComment(line string) bool {
	return line[:1] == "#"
}

func (d *Definition) IsEmptyLine(line string) bool {
	return len(line) == 0
}

func getEntityType(key string) entity.Type {
	if task.Action(key).IsValid() {
		return entity.TypeTask
	}
	return entity.TypeUnknown
}

func splitToKeyValue(str string) (string, string, error) {
	fields := strings.SplitN(str, ":", 2)
	if len(fields) < 2 {
		return "", "", ErrDefInvalidKeyValue
	}

	key, value := fields[0], strings.TrimLeft(fields[1], " \t")
	return key, value, nil
}

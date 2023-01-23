package client

import (
	"errors"
)

var (
	ErrCommandNotParsed error = errors.New("command arguments and flags are not parsed")

	ErrDefUnwantedComment     SyntaxError = "comment is allowed only in the beginning of action definition"
	ErrDefInvalidKeyValue     SyntaxError = "key value pair separated by colon is expected"
	ErrDefUnexpectedAction    SyntaxError = "unexpected action field"
	ErrDefEmptyTarget         SyntaxError = "target name is empty"
	ErrDefUnexpectedEmptyLine SyntaxError = "unexpected empty line"
	ErrDefInvalidField        SyntaxError = "invalid field"
	ErrDefInvalidAction       SyntaxError = "invalid action"
)

type SyntaxError string

func (e SyntaxError) Error() string {
	return string(e)
}

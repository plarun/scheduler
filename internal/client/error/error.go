package error

import "errors"

var (
	ErrInvalidSendEvent error = errors.New("invalid send event")
)

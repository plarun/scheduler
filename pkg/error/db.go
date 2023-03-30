package error

import (
	"fmt"
)

type DatabaseError struct {
	Msg string
}

func NewDatabaseError(msg string) *DatabaseError {
	return &DatabaseError{
		Msg: msg,
	}
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("DB error - %v", e.Msg)
}

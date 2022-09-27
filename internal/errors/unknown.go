package errors

import (
	"fmt"
)

var (
	_ Unknown = (*UnknownError)(nil)
	_ Error   = (*UnknownError)(nil)
)

type Unknown interface {
	error
	IsUnknown()
}

type UnknownError struct {
	*CaosError
}

func ThrowUnknown(parent error, id, message string) error {
	return &UnknownError{CreateCaosError(parent, id, message)}
}

func ThrowUnknownf(parent error, id, format string, a ...interface{}) error {
	return ThrowUnknown(parent, id, fmt.Sprintf(format, a...))
}

func (err *UnknownError) IsUnknown() {}

func IsUnknown(err error) bool {
	_, ok := err.(Unknown)
	return ok
}

func (err *UnknownError) Is(target error) bool {
	t, ok := target.(*UnknownError)
	if !ok {
		return false
	}
	return err.CaosError.Is(t.CaosError)
}

func (err *UnknownError) Unwrap() error {
	return err.CaosError
}

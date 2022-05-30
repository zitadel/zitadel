package errors

import (
	"fmt"
)

var (
	_ Internal = (*InternalError)(nil)
	_ Error    = (*InternalError)(nil)
)

type Internal interface {
	error
	IsInternal()
}

type InternalError struct {
	*CaosError
}

func ThrowInternal(parent error, id, message string) error {
	return &InternalError{CreateCaosError(parent, id, message)}
}

func ThrowInternalf(parent error, id, format string, a ...interface{}) error {
	return ThrowInternal(parent, id, fmt.Sprintf(format, a...))
}

func (err *InternalError) IsInternal() {}

func IsInternal(err error) bool {
	_, ok := err.(Internal)
	return ok
}

func (err *InternalError) Is(target error) bool {
	t, ok := target.(*InternalError)
	if !ok {
		return false
	}
	return err.CaosError.Is(t.CaosError)
}

func (err *InternalError) Unwrap() error {
	return err.CaosError
}

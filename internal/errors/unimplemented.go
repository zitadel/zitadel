package errors

import (
	"fmt"
)

var (
	_ Unimplemented = (*UnimplementedError)(nil)
	_ Error         = (*UnimplementedError)(nil)
)

type Unimplemented interface {
	error
	IsUnimplemented()
}

type UnimplementedError struct {
	*CaosError
}

func ThrowUnimplemented(parent error, id, message string) error {
	return &UnimplementedError{CreateCaosError(parent, id, message)}
}

func ThrowUnimplementedf(parent error, id, format string, a ...interface{}) error {
	return ThrowUnimplemented(parent, id, fmt.Sprintf(format, a...))
}

func (err *UnimplementedError) IsUnimplemented() {}

func IsUnimplemented(err error) bool {
	_, ok := err.(Unimplemented)
	return ok
}

func (err *UnimplementedError) Is(target error) bool {
	t, ok := target.(*UnimplementedError)
	if !ok {
		return false
	}
	return err.CaosError.Is(t.CaosError)
}

func (err *UnimplementedError) Unwrap() error {
	return err.CaosError
}

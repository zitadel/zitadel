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
	return &UnimplementedError{createCaosError(parent, id, message)}
}

func ThrowUnimplementedf(parent error, id, format string, a ...interface{}) error {
	return ThrowUnimplemented(parent, id, fmt.Sprintf(format, a...))
}

func (err *UnimplementedError) IsUnimplemented() {}

func IsUnimplemented(err error) bool {
	_, ok := err.(Unimplemented)
	return ok
}

package errors

import "fmt"

var (
	_ InvalidArgument = (*InvalidArgumentError)(nil)
	_ Error           = (*InvalidArgumentError)(nil)
)

type InvalidArgument interface {
	error
	IsInvalidArgument()
}

type InvalidArgumentError struct {
	*CaosError
}

func ThrowInvalidArgument(parent error, id, message string) error {
	return &InvalidArgumentError{createCaosError(parent, id, message)}
}

func ThrowInvalidArgumentf(parent error, id, format string, a ...interface{}) error {
	return ThrowInvalidArgument(parent, id, fmt.Sprintf(format, a...))
}

func (err *InvalidArgumentError) IsInvalidArgument() {}

func IsErrorInvalidArgument(err error) bool {
	_, ok := err.(InvalidArgument)
	return ok
}

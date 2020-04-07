package errors

import "fmt"

var (
	_ AlreadyExists = (*AlreadyExistsError)(nil)
	_ Error         = (*AlreadyExistsError)(nil)
)

type AlreadyExists interface {
	error
	IsAlreadyExists()
}

type AlreadyExistsError struct {
	*CaosError
}

func ThrowAlreadyExists(parent error, id, message string) error {
	return &AlreadyExistsError{CreateCaosError(parent, id, message)}
}

func ThrowAlreadyExistsf(parent error, id, format string, a ...interface{}) error {
	return &AlreadyExistsError{CreateCaosError(parent, id, fmt.Sprintf(format, a...))}
}

func (err *AlreadyExistsError) IsAlreadyExists() {}

func IsErrorAlreadyExists(err error) bool {
	_, ok := err.(AlreadyExists)
	return ok
}

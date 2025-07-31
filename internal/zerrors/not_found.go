package zerrors

import (
	"errors"
	"fmt"
)

type NotFound interface {
	error
	IsNotFound()
}

type NotFoundError struct {
	*ZitadelError
}

func ThrowNotFound(parent error, id, message string) error {
	return &NotFoundError{CreateZitadelError(parent, id, message)}
}

func ThrowNotFoundf(parent error, id, format string, a ...interface{}) error {
	return ThrowNotFound(parent, id, fmt.Sprintf(format, a...))
}

func (err *NotFoundError) IsNotFound() {}

func IsNotFound(err error) bool {
	var tmp NotFound
	return errors.As(err, &tmp)
}

func (err *NotFoundError) Is(target error) bool {
	t, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return err.ZitadelError.Is(t.ZitadelError)
}

func (err *NotFoundError) Unwrap() error {
	return err.ZitadelError
}

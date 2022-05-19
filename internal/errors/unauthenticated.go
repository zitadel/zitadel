package errors

import (
	"fmt"
)

var (
	_ Unauthenticated = (*UnauthenticatedError)(nil)
	_ Error           = (*UnauthenticatedError)(nil)
)

type Unauthenticated interface {
	error
	IsUnauthenticated()
}

type UnauthenticatedError struct {
	*CaosError
}

func ThrowUnauthenticated(parent error, id, message string) error {
	return &UnauthenticatedError{CreateCaosError(parent, id, message)}
}

func ThrowUnauthenticatedf(parent error, id, format string, a ...interface{}) error {
	return ThrowUnauthenticated(parent, id, fmt.Sprintf(format, a...))
}

func (err *UnauthenticatedError) IsUnauthenticated() {}

func IsUnauthenticated(err error) bool {
	_, ok := err.(Unauthenticated)
	return ok
}

func (err *UnauthenticatedError) Is(target error) bool {
	t, ok := target.(*UnauthenticatedError)
	if !ok {
		return false
	}
	return err.CaosError.Is(t.CaosError)
}

func (err *UnauthenticatedError) Unwrap() error {
	return err.CaosError
}

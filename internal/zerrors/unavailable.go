package zerrors

import (
	"errors"
	"fmt"
)

var (
	_ Unavailable = (*UnavailableError)(nil)
	_ Error       = (*UnavailableError)(nil)
)

type Unavailable interface {
	error
	IsUnavailable()
}

type UnavailableError struct {
	*ZitadelError
}

func ThrowUnavailable(parent error, id, message string) error {
	return &UnavailableError{CreateZitadelError(parent, id, message)}
}

func ThrowUnavailablef(parent error, id, format string, a ...interface{}) error {
	return ThrowUnavailable(parent, id, fmt.Sprintf(format, a...))
}

func (err *UnavailableError) IsUnavailable() {}

func IsUnavailable(err error) bool {
	var tmp Unavailable
	return errors.As(err, &tmp)
}

func (err *UnavailableError) Is(target error) bool {
	t, ok := target.(*UnavailableError)
	if !ok {
		return false
	}
	return err.ZitadelError.Is(t.ZitadelError)
}

func (err *UnavailableError) Unwrap() error {
	return err.ZitadelError
}

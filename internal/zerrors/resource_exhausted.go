package zerrors

import (
	"errors"
	"fmt"
)

var (
	_ ResourceExhausted = (*ResourceExhaustedError)(nil)
	_ Error             = (*ResourceExhaustedError)(nil)
)

type ResourceExhausted interface {
	error
	IsResourceExhausted()
}

type ResourceExhaustedError struct {
	*ZitadelError
}

func ThrowResourceExhausted(parent error, id, message string) error {
	return &ResourceExhaustedError{CreateZitadelError(parent, id, message)}
}

func ThrowResourceExhaustedf(parent error, id, format string, a ...interface{}) error {
	return ThrowResourceExhausted(parent, id, fmt.Sprintf(format, a...))
}

func (err *ResourceExhaustedError) IsResourceExhausted() {}

func IsResourceExhausted(err error) bool {
	var tmp ResourceExhausted
	return errors.As(err, &tmp)
}

func (err *ResourceExhaustedError) Is(target error) bool {
	var tmp *ResourceExhaustedError
	if !errors.As(err, &tmp) {
		return false
	}

	return err.ZitadelError.Is(tmp.ZitadelError)
}

func (err *ResourceExhaustedError) Unwrap() error {
	return err.ZitadelError
}

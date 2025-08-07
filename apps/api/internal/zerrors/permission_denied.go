package zerrors

import (
	"fmt"
)

var (
	_ PermissionDenied = (*PermissionDeniedError)(nil)
	_ Error            = (*PermissionDeniedError)(nil)
)

type PermissionDenied interface {
	error
	IsPermissionDenied()
}

type PermissionDeniedError struct {
	*ZitadelError
}

func ThrowPermissionDenied(parent error, id, message string) error {
	return &PermissionDeniedError{CreateZitadelError(parent, id, message)}
}

func ThrowPermissionDeniedf(parent error, id, format string, a ...interface{}) error {
	return ThrowPermissionDenied(parent, id, fmt.Sprintf(format, a...))
}

func (err *PermissionDeniedError) IsPermissionDenied() {}

func IsPermissionDenied(err error) bool {
	_, ok := err.(PermissionDenied)
	return ok
}

func (err *PermissionDeniedError) Is(target error) bool {
	t, ok := target.(*PermissionDeniedError)
	if !ok {
		return false
	}
	return err.ZitadelError.Is(t.ZitadelError)
}

func (err *PermissionDeniedError) Unwrap() error {
	return err.ZitadelError
}

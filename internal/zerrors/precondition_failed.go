package zerrors

import (
	"errors"
	"fmt"
)

var (
	_ PreconditionFailed = (*PreconditionFailedError)(nil)
	_ Error              = (*PreconditionFailedError)(nil)
)

type PreconditionFailed interface {
	error
	IsPreconditionFailed()
}

type PreconditionFailedError struct {
	*ZitadelError
}

func ThrowPreconditionFailed(parent error, id, message string) error {
	return &PreconditionFailedError{CreateZitadelError(parent, id, message)}
}

func ThrowPreconditionFailedf(parent error, id, format string, a ...interface{}) error {
	return ThrowPreconditionFailed(parent, id, fmt.Sprintf(format, a...))
}

func (err *PreconditionFailedError) IsPreconditionFailed() {}

func IsPreconditionFailed(err error) bool {
	var tmp PreconditionFailed
	return errors.As(err, &tmp)
}

func (err *PreconditionFailedError) Is(target error) bool {
	t, ok := target.(*PreconditionFailedError)
	if !ok {
		return false
	}
	return err.ZitadelError.Is(t.ZitadelError)
}

func (err *PreconditionFailedError) Unwrap() error {
	return err.ZitadelError
}

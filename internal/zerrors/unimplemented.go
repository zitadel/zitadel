package zerrors

import (
	"errors"
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
	*ZitadelError
}

func ThrowUnimplemented(parent error, id, message string) error {
	return &UnimplementedError{CreateZitadelError(parent, id, message)}
}

func ThrowUnimplementedf(parent error, id, format string, a ...interface{}) error {
	return ThrowUnimplemented(parent, id, fmt.Sprintf(format, a...))
}

func (err *UnimplementedError) IsUnimplemented() {}

func IsUnimplemented(err error) bool {
	var tmp Unimplemented
	return errors.As(err, &tmp)
}

func (err *UnimplementedError) Is(target error) bool {
	t, ok := target.(*UnimplementedError)
	if !ok {
		return false
	}
	return err.ZitadelError.Is(t.ZitadelError)
}

func (err *UnimplementedError) Unwrap() error {
	return err.ZitadelError
}

package channels

import "errors"

type ErrCancel struct {
	Err error
}

func (e *ErrCancel) Error() string {
	return e.Err.Error()
}

func NewErrCancel(err error) error {
	return &ErrCancel{
		Err: err,
	}
}

func (e *ErrCancel) Is(target error) bool {
	return errors.As(target, &e)
}

func (e *ErrCancel) Unwrap() error {
	return e.Err
}

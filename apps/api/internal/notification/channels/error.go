package channels

import "errors"

type CancelError struct {
	Err error
}

func (e *CancelError) Error() string {
	return e.Err.Error()
}

func NewCancelError(err error) error {
	return &CancelError{
		Err: err,
	}
}

func (e *CancelError) Is(target error) bool {
	return errors.As(target, &e)
}

func (e *CancelError) Unwrap() error {
	return e.Err
}

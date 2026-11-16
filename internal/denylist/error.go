package denylist

import (
	"fmt"
)

type AddressDeniedError struct {
	deniedBy string
}

func NewAddressDeniedError(deniedBy string) *AddressDeniedError {
	return &AddressDeniedError{deniedBy: deniedBy}
}

func (e *AddressDeniedError) Error() string {
	return fmt.Sprintf("address is denied by '%s'", e.deniedBy)
}

// Is implements the interface used by [errors.Is].
// It only compares the target itself and must not unwrap it:
// callers pass zero values of their own error types as target,
// which are not necessarily safe to unwrap.
func (e *AddressDeniedError) Is(target error) bool {
	addressDeniedErr, ok := target.(*AddressDeniedError)
	if !ok || addressDeniedErr == nil {
		return false
	}
	return e.deniedBy == addressDeniedErr.deniedBy
}

package denylist

import (
	"errors"
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

func (e *AddressDeniedError) Is(target error) bool {
	var addressDeniedErr *AddressDeniedError
	if !errors.As(target, &addressDeniedErr) {
		return false
	}
	return e.deniedBy == addressDeniedErr.deniedBy
}

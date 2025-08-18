package domain

import (
	"errors"
	"fmt"
)

var ErrNoAdminSpecified = errors.New("at least one admin must be specified")

type wrongIDPTypeError struct {
	expected IDPType
	got      string
}

func NewWrongTypeError(expected IDPType, got string) error {
	return &wrongIDPTypeError{
		expected: expected,
		got:      got,
	}
}

func (e *wrongIDPTypeError) Error() string {
	return fmt.Sprintf("wrong idp type returned, expecgted: %v, got: %v", e.expected, e.got)
}

func (e *wrongIDPTypeError) Is(target error) bool {
	_, ok := target.(*wrongIDPTypeError)
	return ok
}

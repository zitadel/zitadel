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

func NewIDPWrongTypeError(expected IDPType, got fmt.Stringer) error {
	return &wrongIDPTypeError{
		expected: expected,
		got:      got.String(),
	}
}

func (e *wrongIDPTypeError) Error() string {
	return fmt.Sprintf("wrong idp type returned, expected: %v, got: %v", e.expected, e.got)
}

func (e *wrongIDPTypeError) Is(target error) bool {
	_, ok := target.(*wrongIDPTypeError)
	return ok
}

type MultipleOrgsUpdatedError struct {
	Msg      string
	Expected int64
	Actual   int64
}

func NewMultipleOrgsUpdatedError(expected, actual int64) error {
	return &MultipleOrgsUpdatedError{
		Expected: expected,
		Actual:   actual,
	}
}

func (err *MultipleOrgsUpdatedError) Error() string {
	return fmt.Sprintf("Message=expecting %d row(s) updated, got %d", err.Expected, err.Actual)
}

type UnexpectedQueryTypeError[T any] struct {
	assertedType T
}

func NewUnexpectedQueryTypeError[T any](assertedType T) error {
	return &UnexpectedQueryTypeError[T]{
		assertedType: assertedType,
	}
}

func (u *UnexpectedQueryTypeError[T]) Error() string {
	return fmt.Sprintf("Message=unexpected query type '%T'", u.assertedType)
}

type UnexpectedTextQueryOperationError[T any] struct {
	assertedType T
}

func NewUnexpectedTextQueryOperationError[T any](assertedType T) error {
	return &UnexpectedTextQueryOperationError[T]{
		assertedType: assertedType,
	}
}

func (u *UnexpectedTextQueryOperationError[T]) Error() string {
	return fmt.Sprintf("Message=unexpected text query operation type '%T'", u.assertedType)
}

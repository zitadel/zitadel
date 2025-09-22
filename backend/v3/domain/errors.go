package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoAdminSpecified = errors.New("at least one admin must be specified")
)

// OrgNotFoundError is used when DB doesn't return a not found (e.g update with no rows updated)
// on organizationlookup but a match is expected
type OrgNotFoundError struct {
	ID string
}

func NewOrgNotFoundError(errID string) error {
	return &OrgNotFoundError{
		ID: errID,
	}
}

func (err *OrgNotFoundError) Error() string {
	return fmt.Sprintf("ID=%s Message=organization not found", err.ID)
}

type MultipleOrgsUpdatedError struct {
	ID       string
	Msg      string
	Expected int64
	Actual   int64
}

func NewMultipleOrgsUpdatedError(id string, expected, actual int64) error {
	return &MultipleOrgsUpdatedError{
		ID:       id,
		Expected: expected,
		Actual:   actual,
	}
}

func (err *MultipleOrgsUpdatedError) Error() string {
	return fmt.Sprintf("ID=%s Message=expecting %d row(s) updated, got %d", err.ID, err.Expected, err.Actual)
}

type OrgNameNotChangedError struct {
	ID string
}

func NewOrgNameNotChangedError(errID string) error {
	return &OrgNameNotChangedError{
		ID: errID,
	}
}

func (err *OrgNameNotChangedError) Error() string {
	return fmt.Sprintf("ID=%s Message=organization name has not changed", err.ID)
}

type UnexpectedQueryTypeError[T any] struct {
	ID           string
	assertedType T
}

func NewUnexpectedQueryTypeError[T any](errID string, assertedType T) error {
	return &UnexpectedQueryTypeError[T]{
		ID:           errID,
		assertedType: assertedType,
	}
}

func (u *UnexpectedQueryTypeError[T]) Error() string {
	return fmt.Sprintf("ID=%s Message=unexpected query type '%T'", u.ID, u.assertedType)
}

type NoQueryCriteriaError struct {
	ID string
}

func NewNoQueryCriteriaError(errID string) error {
	return &NoQueryCriteriaError{
		ID: errID,
	}
}

func (err *NoQueryCriteriaError) Error() string {
	return fmt.Sprintf("ID=%s Message=input query criteria is empty", err.ID)
}

type UnexpectedTextQueryOperationError[T any] struct {
	ID           string
	assertedType T
}

func NewUnexpectedTextQueryOperationError[T any](errID string, assertedType T) error {
	return &UnexpectedTextQueryOperationError[T]{
		ID:           errID,
		assertedType: assertedType,
	}
}

func (u *UnexpectedTextQueryOperationError[T]) Error() string {
	return fmt.Sprintf("ID=%s Message=unexpected text query operation type '%T'", u.ID, u.assertedType)
}

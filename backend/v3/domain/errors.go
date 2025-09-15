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

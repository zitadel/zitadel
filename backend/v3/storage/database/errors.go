package database

import (
	"fmt"
)

// ErrNoRowFound is returned when QueryRow does not find any row.
// It wraps the dialect specific original error to provide more context.
type ErrNoRowFound struct {
	original error
}

func NewNoRowFoundError(original error) error {
	return &ErrNoRowFound{
		original: original,
	}
}

func (e *ErrNoRowFound) Error() string {
	return "no row found"
}

func (e *ErrNoRowFound) Is(target error) bool {
	_, ok := target.(*ErrNoRowFound)
	return ok
}

func (e *ErrNoRowFound) Unwrap() error {
	return e.original
}

// ErrMultipleRowsFound is returned when QueryRow finds multiple rows.
// It wraps the dialect specific original error to provide more context.
type ErrMultipleRowsFound struct {
	original error
	count    int
}

func NewMultipleRowsFoundError(original error, count int) error {
	return &ErrMultipleRowsFound{
		original: original,
		count:    count,
	}
}

func (e *ErrMultipleRowsFound) Error() string {
	return fmt.Sprintf("multiple rows found: %d", e.count)
}

func (e *ErrMultipleRowsFound) Is(target error) bool {
	_, ok := target.(*ErrMultipleRowsFound)
	return ok
}

func (e *ErrMultipleRowsFound) Unwrap() error {
	return e.original
}

type IntegrityType string

const (
	IntegrityTypeCheck   IntegrityType = "check"
	IntegrityTypeUnique  IntegrityType = "unique"
	IntegrityTypeForeign IntegrityType = "foreign"
	IntegrityTypeNotNull IntegrityType = "not null"
	IntegrityTypeUnknown IntegrityType = "unknown"
)

// IntegrityViolation represents a generic integrity violation error.
// It wraps the dialect specific original error to provide more context.
type IntegrityViolation struct {
	integrityType IntegrityType
	table         string
	constraint    string
	original      error
}

func NewIntegrityViolationError(typ IntegrityType, table, constraint string, original error) error {
	return &IntegrityViolation{
		integrityType: typ,
		table:         table,
		constraint:    constraint,
		original:      original,
	}
}

func (e *IntegrityViolation) Error() string {
	return fmt.Sprintf("integrity violation of type %q on %q (constraint: %q): %v", e.integrityType, e.table, e.constraint, e.original)
}

func (e *IntegrityViolation) Is(target error) bool {
	_, ok := target.(*IntegrityViolation)
	return ok
}

// CheckErr is returned when a check constraint fails.
// It wraps the [IntegrityViolation] to provide more context.
// It is used to indicate that a check constraint was violated during an insert or update operation.
type CheckErr struct {
	IntegrityViolation
}

func NewCheckError(table, constraint string, original error) error {
	return &CheckErr{
		IntegrityViolation: IntegrityViolation{
			integrityType: IntegrityTypeCheck,
			table:         table,
			constraint:    constraint,
			original:      original,
		},
	}
}

func (e *CheckErr) Is(target error) bool {
	_, ok := target.(*CheckErr)
	return ok
}

func (e *CheckErr) Unwrap() error {
	return &e.IntegrityViolation
}

// UniqueErr is returned when a unique constraint fails.
// It wraps the [IntegrityViolation] to provide more context.
// It is used to indicate that a unique constraint was violated during an insert or update operation.
type UniqueErr struct {
	IntegrityViolation
}

func NewUniqueError(table, constraint string, original error) error {
	return &UniqueErr{
		IntegrityViolation: IntegrityViolation{
			integrityType: IntegrityTypeUnique,
			table:         table,
			constraint:    constraint,
			original:      original,
		},
	}
}

func (e *UniqueErr) Is(target error) bool {
	_, ok := target.(*UniqueErr)
	return ok
}

func (e *UniqueErr) Unwrap() error {
	return &e.IntegrityViolation
}

// ForeignKeyErr is returned when a foreign key constraint fails.
// It wraps the [IntegrityViolation] to provide more context.
// It is used to indicate that a foreign key constraint was violated during an insert or update operation
type ForeignKeyErr struct {
	IntegrityViolation
}

func NewForeignKeyError(table, constraint string, original error) error {
	return &ForeignKeyErr{
		IntegrityViolation: IntegrityViolation{
			integrityType: IntegrityTypeForeign,
			table:         table,
			constraint:    constraint,
			original:      original,
		},
	}
}

func (e *ForeignKeyErr) Is(target error) bool {
	_, ok := target.(*ForeignKeyErr)
	return ok
}

func (e *ForeignKeyErr) Unwrap() error {
	return &e.IntegrityViolation
}

// NotNullErr is returned when a not null constraint fails.
// It wraps the [IntegrityViolation] to provide more context.
// It is used to indicate that a not null constraint was violated during an insert or update operation.
type NotNullErr struct {
	IntegrityViolation
}

func NewNotNullError(table, constraint string, original error) error {
	return &NotNullErr{
		IntegrityViolation: IntegrityViolation{
			integrityType: IntegrityTypeNotNull,
			table:         table,
			constraint:    constraint,
			original:      original,
		},
	}
}

func (e *NotNullErr) Is(target error) bool {
	_, ok := target.(*NotNullErr)
	return ok
}

func (e *NotNullErr) Unwrap() error {
	return &e.IntegrityViolation
}

// UnknownErr is returned when an unknown error occurs.
// It wraps the dialect specific original error to provide more context.
// It is used to indicate that an error occurred that does not fit into any of the other categories.
type UnknownErr struct {
	original error
}

func NewUnknownError(original error) error {
	return &UnknownErr{
		original: original,
	}
}

func (e *UnknownErr) Error() string {
	return "unknown database error"
}

func (e *UnknownErr) Is(target error) bool {
	_, ok := target.(*UnknownErr)
	return ok
}

func (e *UnknownErr) Unwrap() error {
	return e.original
}

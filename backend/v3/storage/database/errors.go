package database

import (
	"errors"
	"fmt"
)

var (
	ErrNoChanges = errors.New("update must contain a change")
)

type MissingConditionError struct {
	col Column
}

func NewMissingConditionError(col Column) error {
	return &MissingConditionError{
		col: col,
	}
}

func (e *MissingConditionError) Error() string {
	var builder StatementBuilder
	builder.WriteString("missing condition for column")
	if e.col != nil {
		builder.WriteString(" on ")
		e.col.WriteQualified(&builder)
	}
	return builder.String()
}

func (e *MissingConditionError) Is(target error) bool {
	_, ok := target.(*MissingConditionError)
	return ok
}

// NoRowFoundError is returned when QueryRow does not find any row.
// It wraps the dialect specific original error to provide more context.
type NoRowFoundError struct {
	original error
}

func NewNoRowFoundError(original error) error {
	return &NoRowFoundError{
		original: original,
	}
}

func (e *NoRowFoundError) Error() string {
	if e.original != nil {
		return fmt.Sprintf("no row found: %v", e.original)
	}
	return "no row found"
}

func (e *NoRowFoundError) Is(target error) bool {
	_, ok := target.(*NoRowFoundError)
	return ok
}

func (e *NoRowFoundError) Unwrap() error {
	return e.original
}

// MultipleRowsFoundError is returned when QueryRow finds multiple rows.
// It wraps the dialect specific original error to provide more context.
type MultipleRowsFoundError struct {
	original error
}

func NewMultipleRowsFoundError(original error) error {
	return &MultipleRowsFoundError{
		original: original,
	}
}

func (e *MultipleRowsFoundError) Error() string {
	if e.original != nil {
		return fmt.Sprintf("multiple rows found: %v", e.original)
	}
	return "multiple rows found"
}

func (e *MultipleRowsFoundError) Is(target error) bool {
	_, ok := target.(*MultipleRowsFoundError)
	return ok
}

func (e *MultipleRowsFoundError) Unwrap() error {
	return e.original
}

type IntegrityType string

const (
	IntegrityTypeCheck   IntegrityType = "check"
	IntegrityTypeUnique  IntegrityType = "unique"
	IntegrityTypeForeign IntegrityType = "foreign"
	IntegrityTypeNotNull IntegrityType = "not null"
)

// IntegrityViolationError represents a generic integrity violation error.
// It wraps the dialect specific original error to provide more context.
type IntegrityViolationError struct {
	integrityType IntegrityType
	table         string
	constraint    string
	original      error
}

func newIntegrityViolationError(typ IntegrityType, table, constraint string, original error) IntegrityViolationError {
	return IntegrityViolationError{
		integrityType: typ,
		table:         table,
		constraint:    constraint,
		original:      original,
	}
}

func (e *IntegrityViolationError) Error() string {
	if e.original != nil {
		return fmt.Sprintf("integrity violation of type %q on %q (constraint: %q): %v", e.integrityType, e.table, e.constraint, e.original)
	}
	return fmt.Sprintf("integrity violation of type %q on %q (constraint: %q)", e.integrityType, e.table, e.constraint)
}

func (e *IntegrityViolationError) Is(target error) bool {
	_, ok := target.(*IntegrityViolationError)
	return ok
}

func (e *IntegrityViolationError) Unwrap() error {
	return e.original
}

// CheckError is returned when a check constraint fails.
// It wraps the [IntegrityViolationError] to provide more context.
// It is used to indicate that a check constraint was violated during an insert or update operation.
type CheckError struct {
	IntegrityViolationError
}

func NewCheckError(table, constraint string, original error) error {
	return &CheckError{
		IntegrityViolationError: newIntegrityViolationError(IntegrityTypeCheck, table, constraint, original),
	}
}

func (e *CheckError) Is(target error) bool {
	_, ok := target.(*CheckError)
	return ok
}

func (e *CheckError) Unwrap() error {
	return &e.IntegrityViolationError
}

// UniqueError is returned when a unique constraint fails.
// It wraps the [IntegrityViolationError] to provide more context.
// It is used to indicate that a unique constraint was violated during an insert or update operation.
type UniqueError struct {
	IntegrityViolationError
}

func NewUniqueError(table, constraint string, original error) error {
	return &UniqueError{
		IntegrityViolationError: newIntegrityViolationError(IntegrityTypeUnique, table, constraint, original),
	}
}

func (e *UniqueError) Is(target error) bool {
	_, ok := target.(*UniqueError)
	return ok
}

func (e *UniqueError) Unwrap() error {
	return &e.IntegrityViolationError
}

// ForeignKeyError is returned when a foreign key constraint fails.
// It wraps the [IntegrityViolationError] to provide more context.
// It is used to indicate that a foreign key constraint was violated during an insert or update operation
type ForeignKeyError struct {
	IntegrityViolationError
}

func NewForeignKeyError(table, constraint string, original error) error {
	return &ForeignKeyError{
		IntegrityViolationError: newIntegrityViolationError(IntegrityTypeForeign, table, constraint, original),
	}
}

func (e *ForeignKeyError) Is(target error) bool {
	_, ok := target.(*ForeignKeyError)
	return ok
}

func (e *ForeignKeyError) Unwrap() error {
	return &e.IntegrityViolationError
}

// NotNullError is returned when a not null constraint fails.
// It wraps the [IntegrityViolationError] to provide more context.
// It is used to indicate that a not null constraint was violated during an insert or update operation.
type NotNullError struct {
	IntegrityViolationError
}

func NewNotNullError(table, constraint string, original error) error {
	return &NotNullError{
		IntegrityViolationError: newIntegrityViolationError(IntegrityTypeNotNull, table, constraint, original),
	}
}

func (e *NotNullError) Is(target error) bool {
	_, ok := target.(*NotNullError)
	return ok
}

func (e *NotNullError) Unwrap() error {
	return &e.IntegrityViolationError
}

// ScanError is returned when scanning rows into objects failed.
// It wraps the original error to provide more context.
type ScanError struct {
	original error
}

func NewScanError(original error) error {
	return &ScanError{
		original: original,
	}
}

func (e *ScanError) Error() string {
	return fmt.Sprintf("Scan error: %v", e.original)
}

func (e *ScanError) Is(target error) bool {
	_, ok := target.(*ScanError)
	return ok
}

func (e *ScanError) Unwrap() error {
	return e.original
}

// UnknownError is returned when an unknown error occurs.
// It wraps the dialect specific original error to provide more context.
// It is used to indicate that an error occurred that does not fit into any of the other categories.
type UnknownError struct {
	original error
}

func NewUnknownError(original error) error {
	return &UnknownError{
		original: original,
	}
}

func (e *UnknownError) Error() string {
	return fmt.Sprintf("unknown database error: %v", e.original)
}

func (e *UnknownError) Is(target error) bool {
	_, ok := target.(*UnknownError)
	return ok
}

func (e *UnknownError) Unwrap() error {
	return e.original
}

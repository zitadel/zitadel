package domain

import (
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

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

type MultipleObjectsUpdatedError struct {
	Msg      string
	Expected int64
	Actual   int64
}

func NewMultipleObjectsUpdatedError(expected, actual int64) error {
	return &MultipleObjectsUpdatedError{
		Expected: expected,
		Actual:   actual,
	}
}

func (err *MultipleObjectsUpdatedError) Error() string {
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

type PasswordVerificationError struct {
	failedAttempts uint8
}

func NewPasswordVerificationError(failedPassAttempts uint8) error {
	return &PasswordVerificationError{
		failedAttempts: failedPassAttempts,
	}
}

func (e *PasswordVerificationError) Error() string {
	return fmt.Sprintf("Message=failed password attempts (%d)", e.failedAttempts)
}

// handleGetError wraps DB errors coming from Get calls into [zerrors] errors.
//
//   - errorID should be in the format <package short name>-<random id>
//   - objectType should be the string representation of a DB object (e.g. 'user', 'session', 'idp'...)
//
// The function wraps [database.NoRowFoundError] to [zerrors.NotFound] error
// and any other error to [zerrors.InternalError]
func handleGetError(inputErr error, errorID, objectType string) error {
	if inputErr == nil {
		return nil
	}

	if errors.Is(inputErr, &database.NoRowFoundError{}) {
		return zerrors.ThrowNotFoundf(inputErr, errorID, "%s not found", objectType)
	}

	return zerrors.ThrowInternalf(inputErr, errorID, "failed fetching %s", objectType)
}

func handleUpdateError(inputErr error, expectedRowCount, actualRowCount int64, errorID, objectType string) error {
	if inputErr == nil && expectedRowCount == actualRowCount {
		return nil
	}

	if inputErr != nil {
		return zerrors.ThrowInternalf(inputErr, errorID, "failed updating %s", objectType)
	}

	if actualRowCount == 0 {
		return zerrors.ThrowNotFoundf(nil, errorID, "%s not found", objectType)
	}

	if actualRowCount != expectedRowCount {
		return zerrors.ThrowInternal(NewMultipleObjectsUpdatedError(expectedRowCount, actualRowCount), errorID, "unexpected number of rows updated")
	}

	return nil
}

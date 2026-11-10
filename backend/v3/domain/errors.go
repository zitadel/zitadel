package domain

import (
	"errors"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func NewSlug(resource, cause string) zerrors.Slug {
	return zerrors.Slug(fmt.Sprintf("%s.%s", resource, cause))
}

var (
	// SlugInternalError is a general slug for any internal error, where the client doesn't have any influence on the error.
	// It's returned when an unexpected condition in the server occurs, such as a database error, an external service failure, or any other error that is not caused by the client's request.
	// The error details might contain additional information.
	SlugInternalError = NewSlug("zitadel", "internal_error")

	// SlugRequestInvalid is a general slug for any request containing missing, malformed or otherwise invalid parameters.
	// The details typically contain additional information on how to resolve the problem.
	SlugRequestInvalid = NewSlug("request", "invalid")

	// SlugAuthMissingPermission is a general slug for any authorization error where the user is missing a required permission to perform an action.
	// The details typically contain the necessary required permission.
	SlugAuthMissingPermission = NewSlug("auth", "missing_permission")
)

var (
	// ErrIDMissing is an error that can be returned on any resource if the necessary (resource) ID is missing
	ErrIDMissing = func() error {
		return zerrors.CreateZitadelError(zerrors.KindInvalidArgument, nil, string(SlugRequestInvalid), "validation failed: id is required", 1).
			WithDetails(zerrors.ErrorDetailsMap{"id": "required"})
	}

	// ErrInstanceIDMissing is an error that can be returned on any resource if the necessary instance ID is missing
	ErrInstanceIDMissing = func() error {
		return zerrors.CreateZitadelError(zerrors.KindInvalidArgument, nil, string(SlugRequestInvalid), "validation failed: instance_id is required", 1).
			WithDetails(zerrors.ErrorDetailsMap{"instance_id": "required"})
	}

	// ErrInvalidRequest is a general error for any invalid request type, like missing or wrong parameters.
	ErrInvalidRequest = func(message string) error {
		return zerrors.CreateZitadelError(zerrors.KindInvalidArgument, nil, string(SlugRequestInvalid), message, 1)
	}

	// ErrInternal is a general error for any internal error, where the client doesn't have any influence on the error.
	// This error should be used for any error that is not caused by the client's request, but rather by an unexpected condition in the server.
	ErrInternal = func(err error, message string) error {
		return zerrors.CreateZitadelError(zerrors.KindInternal, err, string(SlugInternalError), message, 1)
	}

	// ErrMoreThanOneRowAffected is an error that can be returned when a database operation affects more rows than expected,
	// which could indicate a potential issue with the query or the data integrity.
	ErrMoreThanOneRowAffected = func(message string, rows int64) error {
		return zerrors.CreateZitadelError(zerrors.KindInternal, nil, string(SlugInternalError), message, 1).
			WithDetails(zerrors.ErrorDetailsMap{"rows": rows})
	}

	// ErrAuthMissingPermission is an error that can be returned when a user tries to perform an action that requires a specific permission,
	// but the user does not have that permission.
	ErrAuthMissingPermission = func(err error, message, permission string) error {
		return zerrors.CreateZitadelError(zerrors.KindPermissionDenied, err, string(SlugAuthMissingPermission), message, 1).
			WithDetails(zerrors.ErrorDetailsMap{"permission": permission})
	}
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

func (err *MultipleObjectsUpdatedError) Is(target error) bool {
	_, ok := target.(*MultipleObjectsUpdatedError)
	return ok
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
		return zerrors.CreateZitadelError(zerrors.KindNotFound, inputErr, errorID, fmt.Sprintf("%s not found", objectType), 1)
	}

	return zerrors.CreateZitadelError(zerrors.KindInternal, inputErr, errorID, fmt.Sprintf("failed fetching %s", objectType), 1)
}

func handleUpdateError(inputErr error, expectedRowCount, actualRowCount int64, errorID, objectType string) error {
	if inputErr == nil && expectedRowCount == actualRowCount {
		return nil
	}

	if inputErr != nil {
		return zerrors.CreateZitadelError(zerrors.KindInternal, inputErr, errorID, fmt.Sprintf("failed updating %s", objectType), 1)
	}

	if actualRowCount == 0 {
		return zerrors.CreateZitadelError(zerrors.KindNotFound, nil, errorID, fmt.Sprintf("%s not found", objectType), 1)
	}

	if actualRowCount != expectedRowCount {
		return zerrors.CreateZitadelError(zerrors.KindInternal, NewMultipleObjectsUpdatedError(expectedRowCount, actualRowCount), errorID, "unexpected number of rows updated", 1)
	}

	return nil
}

func (err *PasswordVerificationError) Is(target error) bool {
	_, ok := target.(*PasswordVerificationError)
	return ok
}

type RowsReturnedMismatchError struct {
	Msg      string
	Expected int64
	Actual   int64
}

func NewRowsReturnedMismatchError(expected, actual int64) error {
	return &RowsReturnedMismatchError{
		Expected: expected,
		Actual:   actual,
	}
}

func (err *RowsReturnedMismatchError) Error() string {
	return fmt.Sprintf("Message=expecting %d row(s) returned, got %d", err.Expected, err.Actual)
}

func (err *RowsReturnedMismatchError) Is(target error) bool {
	_, ok := target.(*RowsReturnedMismatchError)
	return ok
}

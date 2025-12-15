package zerrors

import (
	"errors"
	"fmt"
	"log/slog"
)

//go:generate enumer -type=Kind -trimprefix=Kind -json
type Kind int

const (
	// KindCanceled indicates that the operation was canceled, typically by the
	// caller.
	KindCanceled Kind = 1

	// KindUnknown indicates that the operation failed for an unknown reason.
	KindUnknown Kind = 2

	// KindInvalidArgument indicates that client supplied an invalid argument.
	KindInvalidArgument Kind = 3

	// KindDeadlineExceeded indicates that deadline expired before the operation
	// could complete.
	KindDeadlineExceeded Kind = 4

	// KindNotFound indicates that some requested entity (for example, a file or
	// directory) was not found.
	KindNotFound Kind = 5

	// KindAlreadyExists indicates that client attempted to create an entity (for
	// example, a file or directory) that already exists.
	KindAlreadyExists Kind = 6

	// KindPermissionDenied indicates that the caller doesn't have permission to
	// execute the specified operation.
	KindPermissionDenied Kind = 7

	// KindResourceExhausted indicates that some resource has been exhausted. For
	// example, a per-user quota may be exhausted or the entire file system may
	// be full.
	KindResourceExhausted Kind = 8

	// KindPreconditionFailed indicates that the system is not in a state
	// required for the operation's execution.
	KindPreconditionFailed Kind = 9

	// KindAborted indicates that operation was aborted by the system, usually
	// because of a concurrency issue such as a sequencer check failure or
	// transaction abort.
	KindAborted Kind = 10

	// KindOutOfRange indicates that the operation was attempted past the valid
	// range (for example, seeking past end-of-file).
	KindOutOfRange Kind = 11

	// KindUnimplemented indicates that the operation isn't implemented,
	// supported, or enabled in this service.
	KindUnimplemented Kind = 12

	// KindInternal indicates that some invariants expected by the underlying
	// system have been broken. This Kind is reserved for serious errors.
	KindInternal Kind = 13

	// KindUnavailable indicates that the service is currently unavailable. This
	// is usually temporary, so clients can back off and retry idempotent
	// operations.
	KindUnavailable Kind = 14

	// KindDataLoss indicates that the operation has resulted in unrecoverable
	// data loss or corruption.
	KindDataLoss Kind = 15

	// KindUnauthenticated indicates that the request does not have valid
	// authentication credentials for the operation.
	KindUnauthenticated Kind = 16
)

type ZitadelError struct {
	Kind    Kind
	Parent  error
	Message string
	ID      string
}

func ThrowError(parent error, id, message string) error {
	return CreateZitadelError(KindUnknown, parent, id, message)
}

func CreateZitadelError(kind Kind, parent error, id, message string) *ZitadelError {
	return &ZitadelError{
		Kind:    kind,
		Parent:  parent,
		ID:      id,
		Message: message,
	}
}

func (err *ZitadelError) Error() string {
	if err.Parent != nil {
		return fmt.Sprintf("ID=%s Message=%s Parent=(%v)", err.ID, err.Message, err.Parent)
	}
	return fmt.Sprintf("ID=%s Message=%s", err.ID, err.Message)
}

func (err *ZitadelError) Unwrap() error {
	return err.GetParent()
}

func (err *ZitadelError) GetParent() error {
	return err.Parent
}

func (err *ZitadelError) GetMessage() string {
	return err.Message
}

func (err *ZitadelError) SetMessage(msg string) {
	err.Message = msg
}

func (err *ZitadelError) GetID() string {
	return err.ID
}

func (err *ZitadelError) Is(target error) bool {
	t, ok := target.(*ZitadelError)
	if !ok {
		return false
	}
	if t.Kind != err.Kind {
		return false
	}
	if t.ID != "" && t.ID != err.ID {
		return false
	}
	if t.Message != "" && t.Message != err.Message {
		return false
	}
	if t.Parent != nil && !errors.Is(err.Parent, t.Parent) {
		return false
	}

	return true
}

func IsZitadelError(err error) bool {
	zitadelErr := new(ZitadelError)
	return errors.As(err, &zitadelErr)
}

func AsZitadelError(err error) (*ZitadelError, bool) {
	zitadelErr := new(ZitadelError)
	ok := errors.As(err, &zitadelErr)
	return zitadelErr, ok
}

var _ slog.LogValuer = (*ZitadelError)(nil)

func (err *ZitadelError) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("parent", err.Parent.Error()),
		slog.String("message", err.Message),
		slog.String("id", err.ID),
	)
}

package zerrors

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync/atomic"

	"github.com/zitadel/sloggcp"
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

// Because errors are created through singletons, config is global.
var (
	enableReportLocation atomic.Bool
	enableStackTrace     atomic.Bool
	gcpErrorReporting    atomic.Bool
)

// EnableReportLocation enables or disables report locations for created errors.
func EnableReportLocation(enable bool) {
	enableReportLocation.Store(enable)
}

// EnableStackTrace enables or disables stack traces for created errors.
func EnableStackTrace(enable bool) {
	enableStackTrace.Store(enable)
}

// GCPErrorReportingEnabled enables or disables special handling for GCP Error Reporting.
// It must be enabled when using a sloggcp handler to avoid duplicate report locations and stack traces.
func GCPErrorReportingEnabled(enable bool) {
	gcpErrorReporting.Store(enable)
}

type ZitadelError struct {
	Kind    Kind
	Parent  error
	Message string
	ID      string

	// location where the error was created
	reportLocation *sloggcp.ReportLocation
	// stack trace at the point the error was created
	stackTrace    []byte
	hasStackTrace bool
}

func ThrowError(parent error, id, message string) error {
	return newZitadelError(KindUnknown, parent, id, message)
}

func CreateZitadelError(kind Kind, parent error, id, message string) *ZitadelError {
	return newZitadelError(kind, parent, id, message)
}

func newZitadelError(kind Kind, parent error, id, message string) *ZitadelError {
	err := &ZitadelError{
		Kind:    kind,
		Parent:  parent,
		ID:      id,
		Message: message,
	}
	if enableReportLocation.Load() {
		// skip 2: newZitadelError + Create / Throw function.
		err.reportLocation = sloggcp.NewReportLocation(2)
	}
	if enableStackTrace.Load() {
		err.stackTrace = debug.Stack()
		err.hasStackTrace = true
	}
	return err
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
	attributes := make([]slog.Attr, 0, 6)
	attributes = append(attributes,
		slog.String("kind", err.Kind.String()),
		slog.String("message", err.Message),
		slog.String("id", err.ID),
	)
	if err.Parent != nil {
		attributes = append(attributes, slog.Any("parent", err.Parent))
	}
	// if gcp error reporting is enabled, log the error as a group without
	// report location and stack trace, as those are handled by the handler.
	if gcpErrorReporting.Load() {
		return slog.GroupValue(attributes...)
	}

	if err.reportLocation != nil {
		attributes = append(attributes, slog.Any("reportLocation", err.reportLocation))
	}
	if err.hasStackTrace {
		attributes = append(attributes, slog.String("stackTrace", string(err.stackTrace)))
	}
	return slog.GroupValue(attributes...)
}

// ReportLocation implements [sloggcp.ReportLocationError].
func (err *ZitadelError) ReportLocation() *sloggcp.ReportLocation {
	return err.reportLocation
}

// StackTrace implements [sloggcp.StackTraceError].
func (err *ZitadelError) StackTrace() (trace []byte, ok bool) {
	return err.stackTrace, err.hasStackTrace
}

var (
	_ sloggcp.ReportLocationError = (*ZitadelError)(nil)
	_ sloggcp.StackTraceError     = (*ZitadelError)(nil)
)

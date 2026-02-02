// Package logging provides utilities for structured logging with context support.
// It builds on top of slog and slog-context to offer a consistent logging experience
// across different parts of the application by categorizing logs into different streams.
//
// The package uses the global [slog.Default] logger as the base logger,
// which is configured at application startup to set the desired logging level and output format.
// It provides functions to create new loggers for specific streams and to
// add logging capabilities to contexts.
// Streams are defined using the [Stream] enumeration.
// Log context can be created using [NewCtx], and loggers can be retrieved from contexts using [FromCtx].
// Streams are typically initialized at the start of different application components
// (e.g., request handling, event processing) to ensure that all logs generated within those components
// are tagged appropriately.
//
// Example usage:
//
//	// Initialize a context for request handling, typically done in middleware
//	ctx := logging.NewCtx(context.Background(), logging.StreamRequest, slog.String("request_id", "12345"))
//	// Somewhere deeper in the call stack
//	logging.Info(ctx, "Something to log")
//
// This will produce a log entry with the stream set to "request" and include the request ID.
package logging

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"runtime"
	"time"

	slogctx "github.com/veqryn/slog-context"

	"github.com/zitadel/sloggcp"
	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/cmd/build"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// Stream represents a logging stream for categorizing log entries.
// This is a type alias for [instrumentation.Stream] to expose it in this package.
type Stream = instrumentation.Stream

const (
	StreamRuntime      = instrumentation.StreamRuntime      // Application runtime logs.
	StreamReady        = instrumentation.StreamReady        // Readiness and liveness checks.
	StreamRequest      = instrumentation.StreamRequest      // API request handling.
	StreamEventPusher  = instrumentation.StreamEventPusher  // Event pushing to the database (not implemented yet).
	StreamEventHandler = instrumentation.StreamEventHandler // Event handling and processing.
	StreamQueue        = instrumentation.StreamQueue        // Queue operations and job processing.
)

var noop = slog.New(slog.DiscardHandler)

// New creates a new logger with the given stream and additional arguments.
func New(stream Stream, args ...any) *slog.Logger {
	if !instrumentation.IsStreamEnabled(stream) {
		return noop
	}
	args = append(args,
		slog.String("stream", stream.String()),
		slog.String("version", build.Version()),
	)
	return slog.Default().With(args...)
}

// NewCtx creates a new context with a logger for the given stream and additional arguments.
// Use the [FromCtx] or other helpers to retrieve the logger from the context.
// An existing logger in the context will be replaced.
func NewCtx(ctx context.Context, stream Stream, args ...any) context.Context {
	logger := New(stream, args...)
	return ToCtx(ctx, logger)
}

// ToCtx adds the given logger to the context.
// See [slogctx.NewCtx].
func ToCtx(ctx context.Context, logger *slog.Logger) context.Context {
	return slogctx.NewCtx(ctx, logger)
}

// FromCtx retrieves the logger from the context.
// See [slogctx.FromCtx].
func FromCtx(ctx context.Context) *slog.Logger {
	return slogctx.FromCtx(ctx)
}

// With adds the given arguments to the logger in the context.
// See [slogctx.With].
func With(ctx context.Context, args ...any) context.Context {
	return slogctx.With(ctx, args...)
}

// WithGroup adds a group to the logger in the context.
// See [slogctx.WithGroup].
func WithGroup(ctx context.Context, name string) context.Context {
	return slogctx.WithGroup(ctx, name)
}

// Log logs a message with the given level and arguments using the logger from the context.
// See [slogctx.Log].
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	log(ctx, FromCtx(ctx), level, msg, 1, args...)
}

// Debug logs a debug message using the logger from the context.
// See [slogctx.Debug].
func Debug(ctx context.Context, msg string, args ...any) {
	log(ctx, FromCtx(ctx), slog.LevelDebug, msg, 1, args...)
}

// Info logs an info message using the logger from the context.
// See [slogctx.Info].
func Info(ctx context.Context, msg string, args ...any) {
	log(ctx, FromCtx(ctx), slog.LevelInfo, msg, 1, args...)
}

// Warn logs a warning message using the logger from the context.
// See [slogctx.Warn].
func Warn(ctx context.Context, msg string, args ...any) {
	log(ctx, FromCtx(ctx), slog.LevelWarn, msg, 1, args...)
}

// Error logs an error message using the logger from the context.
// See [slogctx.Error].
func Error(ctx context.Context, msg string, args ...any) {
	log(ctx, FromCtx(ctx), slog.LevelError, msg, 1, args...)
}

// WithError adds an error attribute to the logger from the context and returns the new logger.
// If the error is not a [zerrors.ZitadelError], it is wrapped in a generic ZitadelError with kind [zerrors.KindUnknown].
func WithError(ctx context.Context, err error) *ErrorContextLogger {
	var target *zerrors.ZitadelError
	if !errors.As(err, &target) {
		target = zerrors.CreateZitadelError(zerrors.KindUnknown, err, "LOG-Ao5ch", "an unknown error occurred", 1)
	}
	return &ErrorContextLogger{
		ctx:          ctx,
		logger:       slogctx.FromCtx(ctx).With(slogctx.Err(target)),
		canTerminate: true,
	}
}

// OnError adds an error attribute to the logger from the context and returns the new logger
// if the error is not nil. If the error is nil, a no-op logger is returned.
// If the error is not a [zerrors.ZitadelError], it is wrapped in a generic ZitadelError with kind [zerrors.KindUnknown].
func OnError(ctx context.Context, err error) *ErrorContextLogger {
	if err == nil {
		return &ErrorContextLogger{ctx, noop, false}
	}
	var target *zerrors.ZitadelError
	if !errors.As(err, &target) {
		target = zerrors.CreateZitadelError(zerrors.KindUnknown, err, "LOG-ii6Pi", "an unknown error occurred", 1)
	}
	return &ErrorContextLogger{
		ctx:          ctx,
		logger:       slogctx.FromCtx(ctx).With(slogctx.Err(target)),
		canTerminate: true,
	}
}

type ErrorContextLogger struct {
	ctx    context.Context
	logger *slog.Logger
	// canTerminate sets whether Panic/Fatal should actually call panic or os.Exit.
	// False when OnError returned a no-op logger, true in all other cases.
	canTerminate bool
}

func (l *ErrorContextLogger) Debug(msg string, args ...any) {
	log(l.ctx, l.logger, slog.LevelDebug, msg, 1, args...)
}

func (l *ErrorContextLogger) Info(msg string, args ...any) {
	log(l.ctx, l.logger, slog.LevelInfo, msg, 1, args...)
}

func (l *ErrorContextLogger) Warn(msg string, args ...any) {
	log(l.ctx, l.logger, slog.LevelWarn, msg, 1, args...)
}

func (l *ErrorContextLogger) Error(msg string, args ...any) {
	log(l.ctx, l.logger, slog.LevelError, msg, 1, args...)
}

// Panic logs a [sloggcp.LevelAlert] leveled message and panics.
// If the logger was created via [OnError] with a nil error, this method does nothing.
func (l *ErrorContextLogger) Panic(msg string, args ...any) {
	log(l.ctx, l.logger, sloggcp.LevelAlert, msg, 1, args...)
	if l.canTerminate {
		panic(msg)
	}
}

// Fatal logs a [sloggcp.LevelEmergency] leveled message and exits the application with code 1.
// If the logger was created via [OnError] with a nil error, this method does nothing.
func (l *ErrorContextLogger) Fatal(msg string, args ...any) {
	log(l.ctx, l.logger, sloggcp.LevelEmergency, msg, 1, args...)
	if l.canTerminate {
		exit(1)
	}
}

// exit is a variable to allow testing of Fatal without exiting the test process.
var exit = os.Exit

// log is a helper function that logs a message with the given level and arguments using the provided logger.
func log(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, skip int, args ...any) {
	handler := logger.Handler()
	if !handler.Enabled(ctx, level) {
		return
	}
	var pcs [1]uintptr
	if instrumentation.IsAddSourceEnabled() {
		runtime.Callers(skip+2, pcs[:])

	}
	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(ctx, r)
}

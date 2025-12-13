package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	slogctx "github.com/veqryn/slog-context"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func ZitadelErrorToHTTPStatusCode(ctx context.Context, err error) (statusCode int, ok bool) {
	if err == nil {
		return http.StatusOK, true
	}
	statusCode, key, id, lvl := extractError(err)
	msg := key
	msg += " (" + id + ")"
	slogctx.FromCtx(ctx).Log(ctx, lvl, msg, "err", err)
	return statusCode, statusCode != statusUnknown
}

const statusUnknown = 0

func extractError(err error) (statusCode int, msg, id string, lvl slog.Level) {
	zitadelErr := new(zerrors.ZitadelError)
	if ok := errors.As(err, &zitadelErr); !ok {
		return statusUnknown, err.Error(), "", slog.LevelError
	}
	switch {
	case zerrors.IsErrorAlreadyExists(err):
		return http.StatusConflict, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelError
	case zerrors.IsDeadlineExceeded(err):
		return http.StatusGatewayTimeout, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelError
	case zerrors.IsInternal(err):
		return http.StatusInternalServerError, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelError
	case zerrors.IsErrorInvalidArgument(err):
		return http.StatusBadRequest, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelWarn
	case zerrors.IsNotFound(err):
		return http.StatusNotFound, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelWarn
	case zerrors.IsPermissionDenied(err):
		return http.StatusForbidden, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelWarn
	case zerrors.IsPreconditionFailed(err):
		// use the same code as grpc-gateway:
		// https://github.com/grpc-ecosystem/grpc-gateway/blob/9e33e38f15cb7d2f11096366e62ea391a3459ba9/runtime/errors.go#L59
		return http.StatusBadRequest, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelWarn
	case zerrors.IsUnauthenticated(err):
		return http.StatusUnauthorized, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelWarn
	case zerrors.IsUnavailable(err):
		return http.StatusServiceUnavailable, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelError
	case zerrors.IsUnimplemented(err):
		return http.StatusNotImplemented, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelInfo
	case zerrors.IsResourceExhausted(err):
		return http.StatusTooManyRequests, zitadelErr.GetMessage(), zitadelErr.GetID(), slog.LevelError
	default:
		return statusUnknown, err.Error(), "", slog.LevelError
	}
}

func HTTPStatusCodeToZitadelError(parent error, statusCode int, id string, message string) error {
	if statusCode == http.StatusOK {
		return nil
	}
	var errorFunc func(parent error, id, message string) error
	switch statusCode {
	case http.StatusConflict:
		errorFunc = zerrors.ThrowAlreadyExists
	case http.StatusGatewayTimeout:
		errorFunc = zerrors.ThrowDeadlineExceeded
	case http.StatusInternalServerError:
		errorFunc = zerrors.ThrowInternal
	case http.StatusBadRequest:
		errorFunc = zerrors.ThrowInvalidArgument
	case http.StatusNotFound:
		errorFunc = zerrors.ThrowNotFound
	case http.StatusForbidden:
		errorFunc = zerrors.ThrowPermissionDenied
	case http.StatusUnauthorized:
		errorFunc = zerrors.ThrowUnauthenticated
	case http.StatusServiceUnavailable:
		errorFunc = zerrors.ThrowUnavailable
	case http.StatusNotImplemented:
		errorFunc = zerrors.ThrowUnimplemented
	case http.StatusTooManyRequests:
		errorFunc = zerrors.ThrowResourceExhausted
	default:
		errorFunc = zerrors.ThrowUnknown
	}

	return errorFunc(parent, id, message)
}

package http

import (
	"context"
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
	if statusCode == statusUnknown {
		return http.StatusInternalServerError, false
	}
	return statusCode, true
}

const statusUnknown = 0

func extractError(err error) (statusCode int, msg, id string, lvl slog.Level) {
	zitadelErr, ok := zerrors.AsZitadelError(err)
	if !ok {
		return statusUnknown, err.Error(), "", slog.LevelError
	}
	msg, id = zitadelErr.GetMessage(), zitadelErr.GetID()

	switch zitadelErr.Kind {
	case zerrors.KindAlreadyExists:
		statusCode, lvl = http.StatusConflict, slog.LevelError
	case zerrors.KindDeadlineExceeded:
		statusCode, lvl = http.StatusGatewayTimeout, slog.LevelError
	case zerrors.KindInternal:
		statusCode, lvl = http.StatusInternalServerError, slog.LevelError
	case zerrors.KindInvalidArgument:
		statusCode, lvl = http.StatusBadRequest, slog.LevelWarn
	case zerrors.KindNotFound:
		statusCode, lvl = http.StatusNotFound, slog.LevelWarn
	case zerrors.KindPermissionDenied:
		statusCode, lvl = http.StatusForbidden, slog.LevelWarn
	case zerrors.KindPreconditionFailed:
		// use the same code as grpc-gateway:
		// https://github.com/grpc-ecosystem/grpc-gateway/blob/9e33e38f15cb7d2f11096366e62ea391a3459ba9/runtime/errors.go#L59
		statusCode, lvl = http.StatusBadRequest, slog.LevelWarn
	case zerrors.KindUnauthenticated:
		statusCode, lvl = http.StatusUnauthorized, slog.LevelWarn
	case zerrors.KindUnavailable:
		statusCode, lvl = http.StatusServiceUnavailable, slog.LevelError
	case zerrors.KindUnimplemented:
		statusCode, lvl = http.StatusNotImplemented, slog.LevelInfo
	case zerrors.KindResourceExhausted:
		statusCode, lvl = http.StatusTooManyRequests, slog.LevelError
	default:
		statusCode, lvl = statusUnknown, slog.LevelError
	}
	return statusCode, msg, id, lvl
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

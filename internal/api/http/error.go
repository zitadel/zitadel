package http

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func ZitadelErrorToHTTPStatusCode(err error) (statusCode int, ok bool) {
	if err == nil {
		return http.StatusOK, true
	}
	switch {
	case zerrors.IsErrorAlreadyExists(err):
		return http.StatusConflict, true
	case zerrors.IsDeadlineExceeded(err):
		return http.StatusGatewayTimeout, true
	case zerrors.IsInternal(err):
		return http.StatusInternalServerError, true
	case zerrors.IsErrorInvalidArgument(err):
		return http.StatusBadRequest, true
	case zerrors.IsNotFound(err):
		return http.StatusNotFound, true
	case zerrors.IsPermissionDenied(err):
		return http.StatusForbidden, true
	case zerrors.IsPreconditionFailed(err):
		// use the same code as grpc-gateway:
		// https://github.com/grpc-ecosystem/grpc-gateway/blob/9e33e38f15cb7d2f11096366e62ea391a3459ba9/runtime/errors.go#L59
		return http.StatusBadRequest, true
	case zerrors.IsUnauthenticated(err):
		return http.StatusUnauthorized, true
	case zerrors.IsUnavailable(err):
		return http.StatusServiceUnavailable, true
	case zerrors.IsUnimplemented(err):
		return http.StatusNotImplemented, true
	case zerrors.IsResourceExhausted(err):
		return http.StatusTooManyRequests, true
	default:
		return http.StatusInternalServerError, false
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

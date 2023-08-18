package http

import (
	"errors"
	"net/http"

	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func ZitadelErrorToHTTPStatusCode(err error) (statusCode int, ok bool) {
	if err == nil {
		return http.StatusOK, true
	}
	//nolint:errorlint
	switch err.(type) {
	case *caos_errs.AlreadyExistsError:
		return http.StatusConflict, true
	case *caos_errs.DeadlineExceededError:
		return http.StatusGatewayTimeout, true
	case *caos_errs.InternalError:
		return http.StatusInternalServerError, true
	case *caos_errs.InvalidArgumentError:
		return http.StatusBadRequest, true
	case *caos_errs.NotFoundError:
		return http.StatusNotFound, true
	case *caos_errs.PermissionDeniedError:
		return http.StatusForbidden, true
	case *caos_errs.PreconditionFailedError:
		// use the same code as grpc-gateway:
		// https://github.com/grpc-ecosystem/grpc-gateway/blob/9e33e38f15cb7d2f11096366e62ea391a3459ba9/runtime/errors.go#L59
		return http.StatusBadRequest, true
	case *caos_errs.UnauthenticatedError:
		return http.StatusUnauthorized, true
	case *caos_errs.UnavailableError:
		return http.StatusServiceUnavailable, true
	case *caos_errs.UnimplementedError:
		return http.StatusNotImplemented, true
	case *caos_errs.ResourceExhaustedError:
		return http.StatusTooManyRequests, true
	default:
		c := new(caos_errs.CaosError)
		if errors.As(err, &c) {
			return ZitadelErrorToHTTPStatusCode(errors.Unwrap(err))
		}
		return http.StatusInternalServerError, false
	}
}

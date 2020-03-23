package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

func CaosToGRPCError(err error) error {
	if err == nil {
		return nil
	}
	code, msg, ok := extract(err)
	if !ok {
		return status.Convert(err).Err()
	}
	return status.Error(code, msg)
}

func extract(err error) (c codes.Code, msg string, ok bool) {
	switch caosErr := err.(type) {
	case *caos_errs.AlreadyExistsError:
		return codes.AlreadyExists, caosErr.GetMessage(), true
	case *caos_errs.DeadlineExceededError:
		return codes.DeadlineExceeded, caosErr.GetMessage(), true
	case caos_errs.InternalError:
		return codes.Internal, caosErr.GetMessage(), true
	case *caos_errs.InvalidArgumentError:
		return codes.InvalidArgument, caosErr.GetMessage(), true
	case *caos_errs.NotFoundError:
		return codes.NotFound, caosErr.GetMessage(), true
	case *caos_errs.PermissionDeniedError:
		return codes.PermissionDenied, caosErr.GetMessage(), true
	case *caos_errs.PreconditionFailedError:
		return codes.FailedPrecondition, caosErr.GetMessage(), true
	case *caos_errs.UnauthenticatedError:
		return codes.Unauthenticated, caosErr.GetMessage(), true
	case *caos_errs.UnavailableError:
		return codes.Unavailable, caosErr.GetMessage(), true
	case *caos_errs.UnimplementedError:
		return codes.Unimplemented, caosErr.GetMessage(), true
	default:
		return codes.Unknown, err.Error(), false
	}
}

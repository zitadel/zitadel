package errors

import (
	"context"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/pkg/grpc/message"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CaosToGRPCError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	code, key, id, ok := ExtractCaosError(err)
	if !ok {
		return status.Convert(err).Err()
	}
	msg := key
	msg += " (" + id + ")"

	s, err := status.New(code, msg).WithDetails(&message.ErrorDetail{Id: id, Message: key})
	if err != nil {
		logging.Log("GRPC-gIeRw").WithError(err).Debug("unable to add detail")
		return status.New(code, msg).Err()
	}

	return s.Err()
}

func ExtractCaosError(err error) (c codes.Code, msg, id string, ok bool) {
	if err == nil {
		return codes.OK, "", "", false
	}
	switch caosErr := err.(type) {
	case *caos_errs.AlreadyExistsError:
		return codes.AlreadyExists, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.DeadlineExceededError:
		return codes.DeadlineExceeded, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.InternalError:
		return codes.Internal, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.InvalidArgumentError:
		return codes.InvalidArgument, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.NotFoundError:
		return codes.NotFound, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.PermissionDeniedError:
		return codes.PermissionDenied, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.PreconditionFailedError:
		return codes.FailedPrecondition, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.UnauthenticatedError:
		return codes.Unauthenticated, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.UnavailableError:
		return codes.Unavailable, caosErr.GetMessage(), caosErr.GetID(), true
	case *caos_errs.UnimplementedError:
		return codes.Unimplemented, caosErr.GetMessage(), caosErr.GetID(), true
	default:
		return codes.Unknown, err.Error(), "", false
	}
}

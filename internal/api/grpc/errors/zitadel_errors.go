package errors

import (
	"context"

	"github.com/zitadel/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/message"
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
	case *zerrors.AlreadyExistsError:
		return codes.AlreadyExists, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.DeadlineExceededError:
		return codes.DeadlineExceeded, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.InternalError:
		return codes.Internal, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.InvalidArgumentError:
		return codes.InvalidArgument, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.NotFoundError:
		return codes.NotFound, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.PermissionDeniedError:
		return codes.PermissionDenied, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.PreconditionFailedError:
		return codes.FailedPrecondition, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.UnauthenticatedError:
		return codes.Unauthenticated, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.UnavailableError:
		return codes.Unavailable, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.UnimplementedError:
		return codes.Unimplemented, caosErr.GetMessage(), caosErr.GetID(), true
	case *zerrors.ResourceExhaustedError:
		return codes.ResourceExhausted, caosErr.GetMessage(), caosErr.GetID(), true
	default:
		return codes.Unknown, err.Error(), "", false
	}
}

package gerrors

import (
	"github.com/zitadel/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/message"
)

func ZITADELToGRPCError(err error) error {
	if err == nil {
		return nil
	}
	code, key, id, ok := ExtractZITADELError(err)
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

func ExtractZITADELError(err error) (c codes.Code, msg, id string, ok bool) {
	if err == nil {
		return codes.OK, "", "", false
	}
	switch zitadelErr := err.(type) {
	case *zerrors.AlreadyExistsError:
		return codes.AlreadyExists, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.DeadlineExceededError:
		return codes.DeadlineExceeded, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.InternalError:
		return codes.Internal, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.InvalidArgumentError:
		return codes.InvalidArgument, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.NotFoundError:
		return codes.NotFound, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.PermissionDeniedError:
		return codes.PermissionDenied, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.PreconditionFailedError:
		return codes.FailedPrecondition, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.UnauthenticatedError:
		return codes.Unauthenticated, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.UnavailableError:
		return codes.Unavailable, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.UnimplementedError:
		return codes.Unimplemented, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case *zerrors.ResourceExhaustedError:
		return codes.ResourceExhausted, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	default:
		return codes.Unknown, err.Error(), "", false
	}
}

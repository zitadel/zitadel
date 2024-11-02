package gerrors

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
		logging.WithError(err).WithField("logID", "GRPC-gIeRw").Debug("unable to add detail")
		return status.New(code, msg).Err()
	}

	return s.Err()
}

func ExtractZITADELError(err error) (c codes.Code, msg, id string, ok bool) {
	if err == nil {
		return codes.OK, "", "", false
	}
	connErr := new(pgconn.ConnectError)
	if ok := errors.As(err, &connErr); ok {
		return codes.Internal, "db connection error", "", true
	}
	zitadelErr := new(zerrors.ZitadelError)
	if ok := errors.As(err, &zitadelErr); !ok {
		return codes.Unknown, err.Error(), "", false
	}
	switch {
	case zerrors.IsErrorAlreadyExists(err):
		return codes.AlreadyExists, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsDeadlineExceeded(err):
		return codes.DeadlineExceeded, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsInternal(err):
		return codes.Internal, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsErrorInvalidArgument(err):
		return codes.InvalidArgument, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsNotFound(err):
		return codes.NotFound, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsPermissionDenied(err):
		return codes.PermissionDenied, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsPreconditionFailed(err):
		return codes.FailedPrecondition, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsUnauthenticated(err):
		return codes.Unauthenticated, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsUnavailable(err):
		return codes.Unavailable, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsUnimplemented(err):
		return codes.Unimplemented, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	case zerrors.IsResourceExhausted(err):
		return codes.ResourceExhausted, zitadelErr.GetMessage(), zitadelErr.GetID(), true
	default:
		return codes.Unknown, err.Error(), "", false
	}
}

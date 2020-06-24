package grpc

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/i18n"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CaosToGRPCError(ctx context.Context, err error, translator *i18n.Translator) error {
	if err == nil {
		return nil
	}
	code, msg, ok := Extract(err)
	if !ok {
		return status.Convert(err).Err()
	}
	if translator != nil {
		msg = translator.LocalizeFromCtx(ctx, msg, nil)

	}
	return status.Error(code, msg)
}

func Extract(err error) (c codes.Code, msg string, ok bool) {
	switch caosErr := err.(type) {
	case *caos_errs.AlreadyExistsError:
		return codes.AlreadyExists, caosErr.GetMessage(), true
	case *caos_errs.DeadlineExceededError:
		return codes.DeadlineExceeded, caosErr.GetMessage(), true
	case *caos_errs.InternalError:
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

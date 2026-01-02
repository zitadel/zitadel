package gerrors

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/jackc/pgx/v5/pgconn"
	slogctx "github.com/veqryn/slog-context"
	"github.com/zitadel/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"

	commandErrors "github.com/zitadel/zitadel/internal/command/errors"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/message"
)

func ZITADELToGRPCError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	code, key, id, lvl := extractError(err)
	msg := key
	msg += " (" + id + ")"
	slogctx.FromCtx(ctx).Log(ctx, lvl, msg, "err", err, "code", code)

	errorInfo := getErrorInfo(id, key, err)

	s, err := status.New(code, msg).WithDetails(errorInfo)
	if err != nil {
		logging.WithError(err).WithField("logID", "GRPC-gIeRw").Debug("unable to add detail")
		return status.New(code, msg).Err()
	}

	return s.Err()
}

func ZITADELToConnectError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	connectError := new(connect.Error)
	if errors.As(err, &connectError) {
		// Connect error may be returned by other middlewares,
		// so we assume it's a client error and log as warning.
		slogctx.FromCtx(ctx).WarnContext(ctx, connectError.Message(), "err", connectError.Unwrap(), "code", connectError.Code())
		return err
	}
	code, key, id, lvl := extractError(err)
	msg := key
	msg += " (" + id + ")"
	slogctx.FromCtx(ctx).Log(ctx, lvl, msg, "err", err)

	errorInfo := getErrorInfo(id, key, err)

	cErr := connect.NewError(connect.Code(code), errors.New(msg))
	if detail, detailErr := connect.NewErrorDetail(errorInfo.(proto.Message)); detailErr == nil {
		cErr.AddDetail(detail)
	}
	return cErr
}

func ExtractZITADELError(err error) (code codes.Code, msg, id string) {
	if err == nil {
		return codes.OK, "", ""
	}
	code, msg, id, _ = extractError(err)
	return code, msg, id
}

func extractError(err error) (c codes.Code, msg, id string, lvl slog.Level) {
	connErr := new(pgconn.ConnectError)
	if ok := errors.As(err, &connErr); ok {
		return codes.Internal, "db connection error", "", slog.LevelError
	}
	zitadelErr, ok := zerrors.AsZitadelError(err)
	if !ok {
		return codes.Unknown, err.Error(), "", slog.LevelError
	}
	msg, id = zitadelErr.GetMessage(), zitadelErr.GetID()

	switch zitadelErr.Kind {
	case zerrors.KindAlreadyExists:
		c, lvl = codes.AlreadyExists, slog.LevelWarn
	case zerrors.KindDeadlineExceeded:
		c, lvl = codes.DeadlineExceeded, slog.LevelError
	case zerrors.KindInternal:
		c, lvl = codes.Internal, slog.LevelError
	case zerrors.KindInvalidArgument:
		c, lvl = codes.InvalidArgument, slog.LevelWarn
	case zerrors.KindNotFound:
		c, lvl = codes.NotFound, slog.LevelWarn
	case zerrors.KindPermissionDenied:
		c, lvl = codes.PermissionDenied, slog.LevelWarn
	case zerrors.KindPreconditionFailed:
		c, lvl = codes.FailedPrecondition, slog.LevelWarn
	case zerrors.KindUnauthenticated:
		c, lvl = codes.Unauthenticated, slog.LevelWarn
	case zerrors.KindUnavailable:
		c, lvl = codes.Unavailable, slog.LevelError
	case zerrors.KindUnimplemented:
		c, lvl = codes.Unimplemented, slog.LevelInfo
	case zerrors.KindResourceExhausted:
		c, lvl = codes.ResourceExhausted, slog.LevelError
	case zerrors.KindDataLoss:
		c, lvl = codes.DataLoss, slog.LevelError
	case zerrors.KindCanceled:
		c, lvl = codes.Canceled, slog.LevelWarn
	case zerrors.KindAborted:
		c, lvl = codes.Aborted, slog.LevelError
	case zerrors.KindOutOfRange:
		c, lvl = codes.OutOfRange, slog.LevelWarn
	case zerrors.KindUnknown:
		c, lvl = codes.Unknown, slog.LevelError
	default:
		c, lvl = codes.Unknown, slog.LevelError
	}
	return c, msg, id, lvl
}

func getErrorInfo(id, key string, err error) protoadapt.MessageV1 {
	var errorInfo protoadapt.MessageV1

	var wpe *commandErrors.WrongPasswordError
	if err != nil && errors.As(err, &wpe) {
		errorInfo = &message.CredentialsCheckError{Id: id, Message: key, FailedAttempts: wpe.FailedAttempts}
	} else {
		errorInfo = &message.ErrorDetail{Id: id, Message: key}
	}

	return errorInfo
}

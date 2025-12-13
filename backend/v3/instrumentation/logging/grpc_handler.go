package logging

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	slogctx "github.com/veqryn/slog-context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
)

func NewGrpcInterceptor(ignoredMethodSuffixes ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if slices.ContainsFunc(ignoredMethodSuffixes, func(s string) bool {
			return strings.HasSuffix(info.FullMethod, s)
		}) {
			return next(ctx, req)
		}

		logger := instrumentation.Logger()
		ctx = instrumentation.SetGrpcRequestDetails(ctx, info)
		ctx = slogctx.NewCtx(ctx, logger)

		resp, err := next(ctx, req)
		var code codes.Code
		if err != nil {
			code = status.Code(err)
		}
		logger.Log(ctx,
			assertGrpcLevel(code),
			"gRPC served",
			"code", code,
		)

		return resp, err
	}
}

func assertGrpcLevel(code codes.Code) slog.Level {
	switch code {
	case codes.Unimplemented:
		return slog.LevelDebug

	// client errors
	case codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.OutOfRange,
		codes.Unauthenticated:
		return slog.LevelWarn

	// server errors
	case codes.Unknown,
		codes.Aborted,
		codes.Internal,
		codes.Unavailable,
		codes.DataLoss:
		return slog.LevelError

	default: // includes 0 code when no error
		return slog.LevelInfo
	}
}

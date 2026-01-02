package logging

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"connectrpc.com/connect"
	slogctx "github.com/veqryn/slog-context"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
)

func NewConnectInterceptor(next connect.UnaryFunc, ignoredMethodSuffixes ...string) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if slices.ContainsFunc(ignoredMethodSuffixes, func(s string) bool {
			return strings.HasSuffix(req.Spec().Procedure, s)
		}) {
			return next(ctx, req)
		}

		logger := instrumentation.Logger()
		ctx = instrumentation.SetConnectRequestDetails(ctx, req)
		ctx = slogctx.NewCtx(ctx, logger)

		resp, err := next(ctx, req)
		var code connect.Code
		if err != nil {
			code = connect.CodeOf(err)
		}
		logger.Log(ctx,
			assertConnectLevel(code),
			"connect RPC served",
			"code", code,
		)

		return resp, err
	}
}

func assertConnectLevel(code connect.Code) slog.Level {
	switch code {
	case connect.CodeUnimplemented:
		return slog.LevelDebug

	// client errors
	case connect.CodeCanceled,
		connect.CodeInvalidArgument,
		connect.CodeNotFound,
		connect.CodeAlreadyExists,
		connect.CodePermissionDenied,
		connect.CodeResourceExhausted,
		connect.CodeFailedPrecondition,
		connect.CodeOutOfRange,
		connect.CodeUnauthenticated:
		return slog.LevelWarn

	// server errors
	case connect.CodeUnknown,
		connect.CodeAborted,
		connect.CodeDeadlineExceeded,
		connect.CodeInternal,
		connect.CodeUnavailable,
		connect.CodeDataLoss:
		return slog.LevelError

	default: // includes 0 code when no error
		return slog.LevelInfo
	}
}

package logging

import (
	"context"
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
		logger.InfoContext(ctx, "connect RPC request", "code", code)
		return resp, err
	}
}

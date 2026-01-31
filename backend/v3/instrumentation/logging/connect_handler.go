package logging

import (
	"context"
	"slices"
	"strings"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
)

func NewConnectInterceptor(next connect.UnaryFunc, ignoredMethodSuffixes ...string) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if slices.ContainsFunc(ignoredMethodSuffixes, func(s string) bool {
			return strings.HasSuffix(req.Spec().Procedure, s)
		}) {
			return next(ctx, req)
		}

		ctx = NewCtx(ctx, StreamRequest)
		ctx = instrumentation.SetConnectRequestDetails(ctx, req)

		resp, err := next(ctx, req)
		var code connect.Code
		if err != nil {
			code = connect.CodeOf(err)
		}
		Info(ctx, "connect RPC request", "code", code)
		return resp, err
	}
}

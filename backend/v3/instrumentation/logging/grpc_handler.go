package logging

import (
	"context"
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
		logger.InfoContext(ctx, "gRPC request", "code", code)
		return resp, err
	}
}

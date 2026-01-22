package logging

import (
	"context"
	"slices"
	"strings"

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

		ctx = NewCtx(ctx, StreamRequest)
		ctx = instrumentation.SetGrpcRequestDetails(ctx, info)

		resp, err := next(ctx, req)
		var code codes.Code
		if err != nil {
			code = status.Code(err)
		}
		Info(ctx, "gRPC request", "code", code)
		return resp, err
	}
}

package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/api/info"
)

func UnaryActivityClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		activity := info.ActivityInfoFromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx, "zitadel-activity-path", activity.Path)
		ctx = metadata.AppendToOutgoingContext(ctx, "zitadel-activity-request-method", activity.RequestMethod)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

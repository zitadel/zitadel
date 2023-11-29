package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/api/info"
)

func UnaryActivityClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		activityInfo := info.ActivityInfoFromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx, activity.PathKey, activityInfo.Path)
		ctx = metadata.AppendToOutgoingContext(ctx, activity.RequestMethodKey, activityInfo.RequestMethod)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

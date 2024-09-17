package middleware

import (
	"context"
	activity_schemas "github.com/zitadel/zitadel/pkg/streams/activity"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/internal/telemetry/logs/record/activity"
)

func UnaryActivityClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		activityInfo := activity.ActivityInfoFromContext(ctx)
		ctx = metadata.AppendToOutgoingContext(ctx, string(activity_schemas.LogFieldKeyPath), activityInfo.Path)
		ctx = metadata.AppendToOutgoingContext(ctx, string(activity_schemas.LogFieldKeyRequestMethod), activityInfo.RequestMethod)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

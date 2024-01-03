package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/limits"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func LimitsInterceptor(limitsLoader *limits.Loader, ignoreService ...string) grpc.UnaryServerInterceptor {
	for idx, service := range ignoreService {
		if !strings.HasPrefix(service, "/") {
			ignoreService[idx] = "/" + service
		}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
		defer func() { span.EndWithError(err) }()
		for _, service := range ignoreService {
			if strings.HasPrefix(info.FullMethod, service) {
				return handler(ctx, req)
			}
		}
		instance := authz.GetInstance(ctx)
		ctx, l := limitsLoader.Load(interceptorCtx, instance.InstanceID())
		if l.Block != nil && *l.Block {
			return nil, zerrors.ThrowResourceExhausted(nil, "LIMITS-molsj", "Errors.Limits.Instance.Blocked")
		}
		span.End()
		return handler(ctx, req)
	}
}

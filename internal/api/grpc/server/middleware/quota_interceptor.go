package middleware

import (
	"context"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
)

func QuotaExhaustedInterceptor(svc *logstore.Service, ignoreService ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		if !svc.Enabled() {
			return handler(ctx, req)
		}
		interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
		defer func() { span.EndWithError(err) }()

		for _, service := range ignoreService {
			if !strings.HasPrefix(service, "/") {
				service = "/" + service
			}
			if strings.HasPrefix(info.FullMethod, service) {
				return handler(ctx, req)
			}
		}

		instance := authz.GetInstance(ctx)
		remaining := svc.Limit(interceptorCtx, instance.InstanceID())
		if remaining != nil && *remaining == 0 {
			return nil, errors.ThrowResourceExhausted(nil, "QUOTA-vjAy8", "Quota.Access.Exhausted")
		}
		span.End() // TODO: Why?
		return handler(ctx, req)
	}
}

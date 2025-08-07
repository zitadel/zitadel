package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func QuotaExhaustedInterceptor(svc *logstore.Service[*record.AccessLog], ignoreService ...string) grpc.UnaryServerInterceptor {
	for idx, service := range ignoreService {
		if !strings.HasPrefix(service, "/") {
			ignoreService[idx] = "/" + service
		}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		if !svc.Enabled() {
			return handler(ctx, req)
		}
		interceptorCtx, span := tracing.NewServerInterceptorSpan(ctx)
		defer func() { span.EndWithError(err) }()

		// The auth interceptor will ensure that only authorized or public requests are allowed.
		// So if there's no authorization context, we don't need to check for limitation
		// Also, we don't limit calls with system user tokens
		ctxData := authz.GetCtxData(ctx)
		if ctxData.IsZero() || ctxData.SystemMemberships != nil {
			return handler(ctx, req)
		}

		for _, service := range ignoreService {
			if strings.HasPrefix(info.FullMethod, service) {
				return handler(ctx, req)
			}
		}

		instance := authz.GetInstance(ctx)
		remaining := svc.Limit(interceptorCtx, instance.InstanceID())
		if remaining != nil && *remaining == 0 {
			return nil, zerrors.ThrowResourceExhausted(nil, "QUOTA-vjAy8", "Quota.Access.Exhausted")
		}
		span.End()
		return handler(ctx, req)
	}
}

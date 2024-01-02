package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/limits"
	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/logstore/record"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func LimitsInterceptor(logstoreSvc *logstore.Service[*record.AccessLog], limitsLoader *limits.Loader, ignoreService ...string) grpc.UnaryServerInterceptor {
	for idx, service := range ignoreService {
		if !strings.HasPrefix(service, "/") {
			ignoreService[idx] = "/" + service
		}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
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
		// If there is a hard limit on the instance, we block immediately.
		instance := authz.GetInstance(ctx)
		ctx, l := limitsLoader.Load(ctx, instance.InstanceID())
		if l.Block != nil && *l.Block {
			return nil, zerrors.ThrowResourceExhausted(nil, "LIMITS-molsj", "Errors.Limits.Instance.Blocked")
		}
		// If there is no hard limit, we check for a quota
		if !logstoreSvc.Enabled() {
			return handler(ctx, req)
		}
		remaining := logstoreSvc.Limit(interceptorCtx, instance.InstanceID())
		if remaining != nil && *remaining == 0 {
			return nil, zerrors.ThrowResourceExhausted(nil, "QUOTA-vjAy8", "Errors.Quota.Access.Exhausted")
		}
		span.End()
		return handler(ctx, req)
	}
}

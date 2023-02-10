package middleware

import (
	"context"
	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/logstore"
	"google.golang.org/grpc"
	"strings"
)

func QuotaExhaustedInterceptor(svc *logstore.Service, ignoreService ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if !svc.Enabled() {
			return handler(ctx, req)
		}
		for _, service := range ignoreService {
			if !strings.HasPrefix(service, "/") {
				service = "/" + service
			}
			if strings.HasPrefix(info.FullMethod, service) {
				return handler(ctx, req)
			}
		}

		instance := authz.GetInstance(ctx)
		remaining, err := svc.Limit(ctx, instance.InstanceID())
		if err != nil {
			logging.WithError(err).Warn("failed to check whether requests should be limited")
			//nolint:ineffassign
			err = nil
		}

		if remaining != nil && *remaining == 0 {
			return nil, errors.ThrowResourceExhausted(nil, "QUOTA-vjAy8", "Quota.Access.Exhausted")
		}
		return handler(ctx, req)
	}
}

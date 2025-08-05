package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func LimitsInterceptor(ignoreService ...string) grpc.UnaryServerInterceptor {
	for idx, service := range ignoreService {
		if !strings.HasPrefix(service, "/") {
			ignoreService[idx] = "/" + service
		}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		for _, service := range ignoreService {
			if strings.HasPrefix(info.FullMethod, service) {
				return handler(ctx, req)
			}
		}
		instance := authz.GetInstance(ctx)
		if block := instance.Block(); block != nil && *block {
			return nil, zerrors.ThrowResourceExhausted(nil, "LIMITS-molsj", "Errors.Limits.Instance.Blocked")
		}
		return handler(ctx, req)
	}
}

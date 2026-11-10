package connect_middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func LimitsInterceptor(ignoreService ...string) connect.UnaryInterceptorFunc {
	for idx, service := range ignoreService {
		if !strings.HasPrefix(service, "/") {
			ignoreService[idx] = "/" + service
		}
	}

	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (_ connect.AnyResponse, err error) {
			for _, service := range ignoreService {
				if strings.HasPrefix(req.Spec().Procedure, service) {
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
}

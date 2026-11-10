package connect_middleware

import (
	"context"
	"strings"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/internal/api/service"
	_ "github.com/zitadel/zitadel/internal/statik"
)

const (
	unknown = "UNKNOWN"
)

func ServiceHandler() connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			serviceName, _ := serviceAndMethod(req.Spec().Procedure)
			if serviceName != unknown {
				return handler(ctx, req)
			}
			ctx = service.WithService(ctx, serviceName)
			return handler(ctx, req)
		}
	}
}

// serviceAndMethod returns the service and method from a procedure.
func serviceAndMethod(procedure string) (string, string) {
	procedure = strings.TrimPrefix(procedure, "/")
	serviceName, method := unknown, unknown
	if strings.Contains(procedure, "/") {
		long := strings.Split(procedure, "/")[0]
		if strings.Contains(long, ".") {
			split := strings.Split(long, ".")
			serviceName = split[len(split)-1]
		}
	}
	if strings.Contains(procedure, "/") {
		method = strings.Split(procedure, "/")[1]
	}
	return serviceName, method
}

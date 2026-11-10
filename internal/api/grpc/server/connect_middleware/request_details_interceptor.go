package connect_middleware

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// RequestDetailsHandler is a connect interceptor that sets request details in the context
// and adds the ID to the response headers.
// It depends on [CallDurationHandler] to set the request start time in the context.
func RequestDetailsHandler() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			domainCtx := http_util.DomainContext(ctx)
			reqCtx := instrumentation.WithRequestDetails(ctx, domainCtx.InstanceHost, domainCtx.PublicHost)
			id := instrumentation.GetRequestID(reqCtx)

			resp, err := next(reqCtx, req)
			if resp != nil {
				resp.Header().Set(http_util.XRequestID, id.String())
			}

			var target *connect.Error
			if errors.As(err, &target) {
				// adds the request ID to the error response trailer
				target.Meta().Set(http_util.XRequestID, id.String())
			}
			return resp, err
		}
	}
}

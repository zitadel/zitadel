package connect_middleware

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// RequestIDHandler is a connect interceptor that sets a request ID in the context
// and adds it to the response headers.
// It depends on [CallDurationHandler] to set the request start time in the context.
func RequestIDHandler() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			reqCtx, id := instrumentation.NewRequestID(ctx, call.FromContext(ctx))
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

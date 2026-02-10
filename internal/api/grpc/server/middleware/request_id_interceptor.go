package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// RequestIDHandler is a gRPC interceptor that sets a request ID in the context
// and adds it to the response headers.
// It depends on [CallDurationHandler] to set the request start time in the context.
func RequestIDHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqCtx, id := instrumentation.NewRequestID(ctx, call.FromContext(ctx))
		md := metadata.New(map[string]string{http_util.XRequestID: id.String()})
		err := grpc.SetHeader(reqCtx, md)
		logging.OnError(ctx, err).Warn("cannot set request ID to response header")

		return handler(reqCtx, req)
	}
}

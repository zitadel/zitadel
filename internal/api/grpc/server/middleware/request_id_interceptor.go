package middleware

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// RequestIDHandler is a gRPC interceptor that sets a request ID in the context
// and adds it to the response headers.
// It depends on [CallDurationHandler] to set the request start time in the context.
func RequestIDHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		reqCtx, id := instrumentation.NewRequestID(ctx, call.FromContext(ctx))
		resp, err := handler(reqCtx, req)

		md := metadata.New(map[string]string{http_util.XRequestID: id.String()})
		grpc.SetHeader(reqCtx, md)
		return resp, err
	}
}

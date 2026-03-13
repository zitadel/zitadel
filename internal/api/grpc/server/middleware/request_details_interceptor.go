package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// RequestDetailsHandler is a gRPC interceptor that sets a request ID in the context
// and adds the ID to the response headers.
// It depends on [CallDurationHandler] to set the request start time in the context.
func RequestDetailsHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		domainCtx := http_util.DomainContext(ctx)
		reqCtx := instrumentation.WithRequestDetails(ctx, domainCtx.InstanceHost, domainCtx.PublicHost)
		id := instrumentation.GetRequestID(reqCtx)
		md := metadata.New(map[string]string{http_util.XRequestID: id.String()})

		err := grpc.SetHeader(reqCtx, md)
		logging.OnError(reqCtx, err).Warn("cannot set request ID to response header")

		return handler(reqCtx, req)
	}
}

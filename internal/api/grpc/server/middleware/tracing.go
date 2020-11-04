package middleware

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func DefaultTracingServer() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor()
}

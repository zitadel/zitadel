package middleware

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type GRPCMethod string

func DefaultTracingClient() grpc.DialOption {
	return grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor())
}

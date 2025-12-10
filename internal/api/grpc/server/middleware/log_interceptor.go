package middleware

import (
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"google.golang.org/grpc"
)

func LogHandler(ignoredMethodSuffixes ...string) grpc.UnaryServerInterceptor {
	return logging.NewGrpcInterceptor(ignoredMethodSuffixes...)
}

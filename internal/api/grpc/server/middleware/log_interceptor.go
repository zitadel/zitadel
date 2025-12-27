package middleware

import (
	"google.golang.org/grpc"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

func LogHandler(ignoredMethodSuffixes ...string) grpc.UnaryServerInterceptor {
	return logging.NewGrpcInterceptor(ignoredMethodSuffixes...)
}

package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func LogHandler(ignoredMethodSuffixes ...string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		if slices.ContainsFunc(ignoredMethodSuffixes, func(s string) bool {
			return strings.HasSuffix(info.FullMethod, s)
		}) {
			return next(ctx, req)
		}
		start := time.Now()
		ctx = logging.NewCtx(ctx, logging.StreamRequest)
		ctx = instrumentation.SetRequestID(ctx, start)

		resp, err := next(ctx, req)
		var code codes.Code
		if err != nil {
			code = status.Code(err)
		}
		logging.Info(ctx, "request served",
			slog.String("protocol", "grpc"),
			slog.Any("domain", http_util.DomainContext(ctx)),
			slog.String("service", serviceFromRPCMethod(info.FullMethod)),
			slog.String("http_method", http.MethodPost), // gRPC always uses POST
			slog.String("path", info.FullMethod),
			slog.Any("code", code),
			slog.Duration("duration", time.Since(start)),
		)
		return resp, err
	}
}

func serviceFromRPCMethod(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

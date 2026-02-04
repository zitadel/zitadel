package connect_middleware

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"connectrpc.com/connect"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// LogHandler is a connect interceptor that logs the request details
// including protocol, domain, service, HTTP method, path, response code, and duration.
// It depends on [CallDurationHandler] and [RequestIDHandler] to set the request start time and ID in the context.
func LogHandler(ignoredMethodSuffixes ...string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			if slices.ContainsFunc(ignoredMethodSuffixes, func(s string) bool {
				return strings.HasSuffix(req.Spec().Procedure, s)
			}) {
				return next(ctx, req)
			}
			ctx = logging.NewCtx(ctx, logging.StreamRequest)

			resp, err := next(ctx, req)
			var code connect.Code
			if err != nil {
				code = connect.CodeOf(err)
			}
			spec := req.Spec()
			logging.Info(ctx, "request served",
				slog.String("protocol", "connect"),
				slog.Any("domain", http_util.DomainContext(ctx)),
				slog.String("service", serviceFromRPCMethod(spec.Procedure)),
				slog.String("http_method", req.HTTPMethod()),
				slog.String("path", spec.Procedure),
				slog.Any("code", code),
				slog.Duration("duration", call.Took(ctx)),
			)
			return resp, err
		}
	}
}

func serviceFromRPCMethod(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) >= 2 {
		return parts[1]
	}
	return "unknown"
}

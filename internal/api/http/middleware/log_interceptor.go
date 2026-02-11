package middleware

import (
	"log/slog"
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// LogHandler is an HTTP middleware that logs the request details
// including protocol, domain, service, HTTP method, path, response code, and duration.
// It depends on [CallDurationHandler] and [RequestIDHandler] to set the request start time and ID in the context.
func LogHandler(service string, ignoredPrefix ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		filter := instrumentation.RequestFilter(ignoredPrefix...)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !filter(r) {
				next.ServeHTTP(w, r)
				return
			}
			ctx := logging.NewCtx(r.Context(), logging.StreamRequest)
			sw := newStatusWriter(w)

			next.ServeHTTP(sw, r.WithContext(ctx))

			logging.Info(ctx, "request served",
				slog.String("protocol", "http"),
				slog.Any("domain", http_util.DomainContext(ctx)),
				slog.String("service", service),
				slog.String("http_method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", sw.status),
				slog.Duration("duration", call.Took(ctx)),
			)
		})
	}
}

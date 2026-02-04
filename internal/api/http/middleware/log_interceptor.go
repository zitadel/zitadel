package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func LogHandler(service string, ignoredPrefix ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		filter := instrumentation.RequestFilter(ignoredPrefix...)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !filter(r) {
				next.ServeHTTP(w, r)
				return
			}
			start := time.Now()
			ctx := logging.NewCtx(r.Context(), logging.StreamRequest)
			ctx = instrumentation.SetRequestID(ctx, start)
			sw := newStatusWriter(w)

			next.ServeHTTP(sw, r.WithContext(ctx))

			logging.Info(ctx, "request served",
				slog.String("protocol", "http"),
				slog.Any("domain", http_util.DomainContext(ctx)),
				slog.String("service", service),
				slog.String("http_method", r.Method), // gRPC always uses POST
				slog.String("path", r.URL.Path),
				slog.Int("status", sw.status),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

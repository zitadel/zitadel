package logging

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func NewHandler(next http.Handler, service string, ignoredPrefix ...string) http.Handler {
	filter := instrumentation.RequestFilter(ignoredPrefix...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !filter(r) {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		ctx := NewCtx(r.Context(), StreamRequest)
		ctx = instrumentation.SetRequestID(ctx, start)
		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r.WithContext(ctx))

		Info(ctx, "request served",
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

// statusWriter is a [http.ResponseWriter] that captures the status code for logging.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) Write(p []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	return w.ResponseWriter.Write(p)
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

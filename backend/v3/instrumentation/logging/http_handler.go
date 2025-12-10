package logging

import (
	"log/slog"
	"net/http"
	"time"

	slogctx "github.com/veqryn/slog-context"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
)

func NewHandler(next http.Handler, service string, ignoredPrefix ...string) http.Handler {
	filter := instrumentation.RequestFilter(ignoredPrefix...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !filter(r) {
			next.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		logger := instrumentation.Logger()
		ctx := instrumentation.SetHttpRequestDetails(r.Context(), service, r)
		ctx = slogctx.NewCtx(ctx, logger)
		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r.WithContext(ctx))

		logger.Log(ctx,
			sw.logLevel(),
			"http request served",
			"status", sw.status,
			"duration", time.Since(start),
		)
	})
}

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

func (w statusWriter) logLevel() slog.Level {
	if w.status < 200 {
		return slog.LevelDebug
	}
	if w.status < 400 {
		return slog.LevelInfo
	}
	if w.status < 500 {
		return slog.LevelWarn
	}
	return slog.LevelError
}

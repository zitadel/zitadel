package logging

import (
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
)

func NewHandler(next http.Handler, service string, ignoredPrefix ...string) http.Handler {
	filter := instrumentation.RequestFilter(ignoredPrefix...)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !filter(r) {
			next.ServeHTTP(w, r)
			return
		}
		ctx := NewCtx(r.Context(), StreamRequest)
		ctx = instrumentation.SetHttpRequestDetails(ctx, service, r)
		sw := &statusWriter{ResponseWriter: w}

		next.ServeHTTP(sw, r.WithContext(ctx))

		Info(ctx, "http request", "status", sw.status)
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

package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/activity"
	"github.com/zitadel/zitadel/internal/api/info"
)

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *statusWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}

func ActivityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := info.ActivityInfoFromContext(r.Context()).SetPath(r.URL.Path).SetRequestMethod(r.Method).IntoContext(r.Context())
		ctx = activity.CreateStorageInfoContext(ctx)
		sw := &statusWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(sw, r.WithContext(ctx))
		activity.TriggerHTTP(ctx, sw.statusCode)
	})
}

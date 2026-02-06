package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/internal/api/call"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// RequestIDHandler is a HTTP middleware that sets a request ID in the context
// and adds it to the response headers.
// It depends on [CallDurationHandler] to set the request start time in the context.
func RequestIDHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, id := instrumentation.NewRequestID(r.Context(), call.FromContext(r.Context()))
			w.Header().Set(http_util.XRequestID, id.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	http_util "github.com/zitadel/zitadel/internal/api/http"
)

// RequestDetailsHandler is a HTTP middleware that sets a request ID in the context
// and adds the ID to the response headers.
// It depends on [CallDurationHandler] to set the request start time in the context.
func RequestDetailsHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			domainCtx := http_util.DomainContext(r.Context())
			ctx := instrumentation.WithRequestDetails(r.Context(), domainCtx.InstanceHost, domainCtx.PublicHost)
			id := instrumentation.GetRequestID(ctx)
			w.Header().Set(http_util.XRequestID, id.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

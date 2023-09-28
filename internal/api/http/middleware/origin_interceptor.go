package middleware

import (
	"net/http"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		origin := r.Header.Get(http_util.Origin)
		if origin != "" {
			ctx = http_util.WithOrigin(ctx, origin)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

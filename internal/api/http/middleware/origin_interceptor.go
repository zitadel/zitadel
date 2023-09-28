package middleware

import (
	"net/http"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get(http_util.Origin)
		if !http_util.IsOrigin(origin) {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(http_util.WithOrigin(r.Context(), origin)))
	})
}

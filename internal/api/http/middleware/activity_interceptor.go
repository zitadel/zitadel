package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/v2/internal/api/info"
)

func ActivityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := info.ActivityInfoFromContext(r.Context()).SetPath(r.URL.Path).SetRequestMethod(r.Method).IntoContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

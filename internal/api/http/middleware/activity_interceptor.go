package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/api/info"
)

func ActivityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(info.RequestMethodIntoContext(r.Method)(info.HTTPPathIntoContext(r.URL.Path)(r.Context()))))
	})
}

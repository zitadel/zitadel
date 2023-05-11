package middleware

import (
	"net/http"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

func RobotsTagHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(http_utils.XRobotsTag, "none")
		next.ServeHTTP(w, r)
	})
}

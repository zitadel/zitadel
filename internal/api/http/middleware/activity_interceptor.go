package middleware

import (
	"net/http"
	"strings"

	"github.com/zitadel/zitadel/internal/api/info"
)

func ActivityHandler(handlerPrefixes []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			activityInfo := info.ActivityInfoFromContext(ctx)
			hasPrefix := false
			// only add path to context if handler is called
			for _, prefix := range handlerPrefixes {
				if strings.HasPrefix(r.URL.Path, prefix) {
					activityInfo.SetPath(r.URL.Path)
					hasPrefix = true
					break
				}
			}
			// last call is with grpc method as path
			if !hasPrefix {
				activityInfo.SetMethod(r.URL.Path)
			}
			ctx = activityInfo.SetRequestMethod(r.Method).IntoContext(ctx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

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
			hasPrefix := false
			// only add path to context if handler is called
			for _, prefix := range handlerPrefixes {
				if strings.HasPrefix(r.URL.Path, prefix) {
					ctx = info.HTTPPathIntoContext(r.URL.Path)(ctx)
					hasPrefix = true
					break
				}
			}
			// last call is with grpc method as path
			if !hasPrefix {
				ctx = info.RPCMethodIntoContext(r.URL.Path)(ctx)
			}
			next.ServeHTTP(w, r.WithContext(info.RequestMethodIntoContext(r.Method)(ctx)))
		})
	}
}

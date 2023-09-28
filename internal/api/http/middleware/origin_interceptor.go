package middleware

import (
	"net/http"
	"strconv"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		forwardedPort := r.Header.Get("x-forwarded-port")
		forwardedHost := r.Header.Get("x-forwarded-host")
		forwardedProto := r.Header.Get("x-forwarded-proto")
		forwardedProtoBool := forwardedProto == "https"
		forwardedPortInt, err := strconv.ParseUint(forwardedPort, 10, 16)
		var origin string
		if err == nil {
			origin = http_util.BuildHTTP(forwardedHost, uint16(forwardedPortInt), forwardedProtoBool)
		}
		if origin == "" {
			origin = r.Header.Get(http_util.Origin)
		}
		if http_util.IsOrigin(origin) {
			ctx = http_util.WithOrigin(ctx, origin)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

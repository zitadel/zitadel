package middleware

import (
	"net/http"

	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.String()
		origin, err := http_util.GetOriginFromURLString(url)
		logging.OnError(err).Debugf("can't get origin from url: %v", url)
		ctx := r.Context()
		if origin != "" {
			ctx = http_util.WithOrigin(ctx, origin)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

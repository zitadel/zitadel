package middleware

import (
	"fmt"
	"net/http"

	"github.com/muhlemmer/httpforwarded"
	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := composeOrigin(r)
		if !http_util.IsOrigin(origin) {
			logging.Debugf("extracted origin is not valid: %s", origin)
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(http_util.WithComposedOrigin(r.Context(), origin)))
	})
}

func composeOrigin(r *http.Request) string {
	var proto, host string
	fwd, fwdErr := httpforwarded.ParseFromRequest(r)
	if fwdErr == nil {
		proto = oldestForwardedValue(fwd, "proto")
		host = oldestForwardedValue(fwd, "host")
	}
	if proto == "" {
		proto = r.Header.Get("X-Forwarded-Proto")
	}
	if host == "" {
		host = r.Header.Get("X-Forwarded-Host")
	}
	if proto == "" {
		if r.TLS == nil {
			proto = "http"
		} else {
			proto = "https"
		}
	}
	if host == "" {
		host = r.Host
	}
	return fmt.Sprintf("%s://%s", proto, host)
}

func oldestForwardedValue(forwarded map[string][]string, key string) string {
	if forwarded == nil {
		return ""
	}
	values := forwarded[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

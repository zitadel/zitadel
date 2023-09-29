package middleware

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/muhlemmer/httpforwarded"
	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin, err := buildOrigin(r)
		if err != nil || !http_util.IsOrigin(origin) {
			logging.OnError(err).Debug("failed to build origin from request")
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(http_util.WithOrigin(r.Context(), origin)))
	})
}

func buildOrigin(r *http.Request) (string, error) {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	origin := fmt.Sprintf("%s://%s", scheme, r.Host)
	fwd, err := httpforwarded.ParseFromRequest(r)
	if err != nil {
		return origin, err
	}
	var fwdProto, fwdHost, fwdPort string
	if fwdProto = mostRecentValue(fwd, "proto"); fwdProto == "" {
		return origin, nil
	}
	if fwdHost = mostRecentValue(fwd, "host"); fwdHost == "" {
		return origin, nil
	}
	fwdPort, foundFwdFor := extractPort(mostRecentValue(fwd, "for"))
	if !foundFwdFor {
		return origin, nil
	}
	o := fmt.Sprintf("%s://%s", fwdProto, fwdHost)
	if fwdPort != "" {
		o += ":" + fwdPort
	}
	return o, nil
}

func extractPort(raw string) (string, bool) {
	if u, ok := parseURL(raw); ok {
		return u.Port(), ok
	}
	return "", false
}

func parseURL(raw string) (*url.URL, bool) {
	if raw == "" {
		return nil, false
	}
	u, err := url.Parse(raw)
	return u, err == nil
}

func mostRecentValue(forwarded map[string][]string, key string) string {
	if forwarded == nil {
		return ""
	}
	values := forwarded[key]
	if len(values) == 0 {
		return ""
	}
	return values[len(values)-1]
}

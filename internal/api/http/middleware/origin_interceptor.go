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
		origin := buildOrigin(r)
		if !http_util.IsOrigin(origin) {
			logging.Debugf("extracted origin is not valid: %s", origin)
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(http_util.WithOrigin(r.Context(), origin)))
	})
}

func buildOrigin(r *http.Request) string {
	if origin, err := originFromForwardedHeader(r); err != nil {
		logging.OnError(err).Debug("failed to build origin from forwarded header, trying x-forwarded-* headers")
	} else {
		return origin
	}
	if origin, err := originFromXForwardedHeaders(r); err != nil {
		logging.OnError(err).Debug("failed to build origin from x-forwarded-* headers, using host header")
	} else {
		return origin
	}
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func originFromForwardedHeader(r *http.Request) (string, error) {
	fwd, err := httpforwarded.ParseFromRequest(r)
	if err != nil {
		return "", err
	}
	var fwdProto, fwdHost, fwdPort string
	if fwdProto = mostRecentValue(fwd, "proto"); fwdProto == "" {
		return "", fmt.Errorf("no proto in forwarded header")
	}
	if fwdHost = mostRecentValue(fwd, "host"); fwdHost == "" {
		return "", fmt.Errorf("no host in forwarded header")
	}
	fwdPort, foundFwdFor := extractPort(mostRecentValue(fwd, "for"))
	if !foundFwdFor {
		return "", fmt.Errorf("no for in forwarded header")
	}
	o := fmt.Sprintf("%s://%s", fwdProto, fwdHost)
	if fwdPort != "" {
		o += ":" + fwdPort
	}
	return o, nil
}

func originFromXForwardedHeaders(r *http.Request) (string, error) {
	scheme := r.Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		return "", fmt.Errorf("no X-Forwarded-Proto header")
	}
	host := r.Header.Get("X-Forwarded-Host")
	if host == "" {
		return "", fmt.Errorf("no X-Forwarded-Host header")
	}
	return fmt.Sprintf("%s://%s", scheme, host), nil
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

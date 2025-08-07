package middleware

import (
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"github.com/muhlemmer/httpforwarded"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func WithOrigin(enforceHttps bool, http1Header, http2Header string, instanceHostHeaders, publicDomainHeaders []string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := composeDomainContext(
				r,
				enforceHttps,
				// to make sure we don't break existing configurations we append the existing checked headers as well
				slices.Compact(append(instanceHostHeaders, http1Header, http2Header, http_util.Forwarded, http_util.ZitadelForwarded, http_util.ForwardedFor, http_util.ForwardedHost, http_util.ForwardedProto)),
				publicDomainHeaders,
			)
			next.ServeHTTP(w, r.WithContext(http_util.WithDomainContext(r.Context(), origin)))
		})
	}
}

func composeDomainContext(r *http.Request, enforceHttps bool, instanceDomainHeaders, publicDomainHeaders []string) *http_util.DomainCtx {
	instanceHost, instanceProto := hostFromRequest(r, instanceDomainHeaders)
	publicHost, publicProto := hostFromRequest(r, publicDomainHeaders)
	if instanceHost == "" {
		instanceHost = r.Host
	}
	return &http_util.DomainCtx{
		InstanceHost: instanceHost,
		Protocol:     protocolFromRequest(instanceProto, publicProto, enforceHttps),
		PublicHost:   publicHost,
	}
}

func protocolFromRequest(instanceProto, publicProto string, enforceHttps bool) string {
	if enforceHttps {
		return "https"
	}
	if publicProto != "" {
		return publicProto
	}
	if instanceProto != "" {
		return instanceProto
	}
	return "http"
}

func hostFromRequest(r *http.Request, headers []string) (host, proto string) {
	var hostFromHeader, protoFromHeader string
	for _, header := range headers {
		switch http.CanonicalHeaderKey(header) {
		case http.CanonicalHeaderKey(http_util.Forwarded),
			http.CanonicalHeaderKey(http_util.ForwardedFor),
			http.CanonicalHeaderKey(http_util.ZitadelForwarded):
			hostFromHeader, protoFromHeader = hostFromForwarded(r.Header.Values(header))
		case http.CanonicalHeaderKey(http_util.ForwardedHost):
			hostFromHeader = r.Header.Get(header)
		case http.CanonicalHeaderKey(http_util.ForwardedProto):
			protoFromHeader = r.Header.Get(header)
		default:
			hostFromHeader = r.Header.Get(header)
		}
		if host == "" {
			host = hostFromHeader
		}
		if proto == "" && (protoFromHeader == "http" || protoFromHeader == "https") {
			proto = protoFromHeader
		}
	}
	return host, proto
}

func hostFromForwarded(values []string) (string, string) {
	fwd, fwdErr := httpforwarded.Parse(values)
	if fwdErr == nil {
		return oldestForwardedValue(fwd, "host"), oldestForwardedValue(fwd, "proto")
	}
	return "", ""
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

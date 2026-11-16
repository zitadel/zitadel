package middleware

import (
	"errors"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"github.com/muhlemmer/httpforwarded"
	"github.com/zitadel/logging"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func WithOrigin(enforceHttps bool, http1Header, http2Header string, instanceHostHeaders, publicDomainHeaders []string) mux.MiddlewareFunc {
	// to make sure we don't break existing configurations, we append the existing checked headers as well
	instanceDomainHeaders := slices.Compact(append(slices.Clone(instanceHostHeaders), http1Header, http2Header, http_util.Forwarded, http_util.ZitadelForwarded, http_util.ForwardedFor, http_util.ForwardedHost, http_util.ForwardedProto))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin, err := composeDomainContext(r, enforceHttps, instanceDomainHeaders, publicDomainHeaders)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r.WithContext(http_util.WithDomainContext(r.Context(), origin)))
		})
	}
}

// composeDomainContext only checks the hosts structurally.
// Whether the hosts belong to an instance and the public host is trusted is checked by the instance interceptors.
func composeDomainContext(r *http.Request, enforceHttps bool, instanceDomainHeaders, publicDomainHeaders []string) (*http_util.DomainCtx, error) {
	// Invalid instance host headers are ignored and fall back to the request host.
	instanceHost, instanceProto, _ := hostFromRequest(r, instanceDomainHeaders)
	// An invalid public host must not be ignored: an empty public host falls back to the instance host
	// and skips the trusted domain check, which would allow building links (e.g. password reset) with an untrusted host.
	publicHost, publicProto, err := hostFromRequest(r, publicDomainHeaders)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "HTTP-Rk7pX", "invalid public host")
	}
	if instanceHost == "" {
		instanceHost = r.Host
	}
	return http_util.NewDomainCtx(instanceHost, publicHost, protocolFromRequest(instanceProto, publicProto, enforceHttps)), nil
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

// hostFromRequest returns the first valid host and the first valid protocol found in the headers,
// in the configured order: the first header that yields a valid host wins (first match wins), and a
// lower-ranked header is only consulted for the host when the earlier ones carry none. The returned
// error belongs to the header that provided the host (e.g. ambiguous multiple values), or, if no
// header provides a host, to the first header that carried an invalid value.
// Callers decide whether that error is ignored (instance host) or rejects the request (public host).
func hostFromRequest(r *http.Request, headers []string) (host, proto string, err error) {
	for _, header := range headers {
		values := r.Header.Values(header)
		headerHost, headerProto, headerErr := hostFromHeader(header, values)
		if proto == "" {
			proto = headerProto
		}
		if host != "" {
			// the first valid host already won; later headers only contribute a protocol
			continue
		}
		if headerHost != "" {
			// the winning host header: keep its own error, e.g. ambiguous multiple values
			host, err = headerHost, headerErr
			continue
		}
		if headerErr != nil {
			logging.WithFields("header", header, "values", values).WithError(headerErr).Warning("invalid host in header, ignoring header")
			if err == nil {
				err = headerErr
			}
		}
	}
	return host, proto, err
}

// hostFromHeader returns the validated host and the protocol ("http" or "https") of a header.
// A plain host header (e.g. X-Forwarded-Host or a custom header) with more than one value is
// ambiguous: the first value is returned together with an error, so it can still be used for the
// instance host (e.g. chained proxies), but rejects a public host, which must be unambiguous since
// it is used to build links.
// A structured Forwarded / Zitadel-Forwarded header instead carries one element per proxy hop,
// so the oldest (client-facing) host and proto are used and further hops are not treated as ambiguous.
func hostFromHeader(header string, values []string) (host, proto string, err error) {
	if len(values) == 0 {
		return "", "", nil
	}
	var rawHost string
	switch http.CanonicalHeaderKey(header) {
	case http.CanonicalHeaderKey(http_util.ForwardedFor):
		// X-Forwarded-For only contains client addresses
		return "", "", nil
	case http.CanonicalHeaderKey(http_util.Forwarded), http.CanonicalHeaderKey(http_util.ZitadelForwarded):
		forwarded, err := httpforwarded.Parse(values)
		if err != nil {
			return "", "", zerrors.ThrowInvalidArgument(err, "HTTP-Jc3nW", "invalid forwarded header")
		}
		rawHost, proto = oldestForwardedValue(forwarded, "host"), oldestForwardedValue(forwarded, "proto")
	case http.CanonicalHeaderKey(http_util.ForwardedProto):
		proto = values[0]
	default:
		rawHost = values[0]
		if len(values) > 1 {
			err = zerrors.ThrowInvalidArgument(nil, "HTTP-Tg5yB", "multiple host values")
		}
	}
	if proto != "http" && proto != "https" {
		proto = ""
	}
	if rawHost == "" {
		return "", proto, err
	}
	host, hostErr := http_util.ValidateHost(rawHost)
	return host, proto, errors.Join(err, hostErr)
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

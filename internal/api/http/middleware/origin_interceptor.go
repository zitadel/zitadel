package middleware

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/muhlemmer/httpforwarded"
	"google.golang.org/grpc/metadata"

	http_util "github.com/zitadel/zitadel/internal/api/http"
)

func OriginHandlerFunc(externalSecure bool, http1HostHeaderOverwrite, http2HostHeaderOverwrite string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin, err := composeOrigin(r, externalSecure, http1HostHeaderOverwrite, http2HostHeaderOverwrite)
			if err != nil {
				// TODO: translate?
				http.Error(w, fmt.Sprintf("invalid origin: %v", err), http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r.WithContext(http_util.WithRequestOrigin(r.Context(), origin)))
		})
	}
}

func composeOrigin(r *http.Request, externalSecure bool, http1HostHeaderOverwrite, http2HostHeaderOverwrite string) (origin http_util.RequestOrigin, err error) {
	http2HostHeader := ":authority"
	if http2HostHeaderOverwrite != "" {
		http2HostHeader = http2HostHeaderOverwrite
	}
	origin.Host = r.Header.Get(http2HostHeader)
	if origin.Host == "" && http1HostHeaderOverwrite != "" {
		origin.Host = r.Header.Get(http1HostHeaderOverwrite)
	}
	fwd, fwdErr := httpforwarded.ParseFromRequest(r)
	if fwdErr == nil {
		/* TODO: It doesn't make sense to read the proto directive from the Forwarded header as long as the saml and oidc libraries don't do so too.
		if origin.Scheme == "" {
			origin.Scheme = oldestForwardedValue(fwd, "proto")
		}
		*/
		if origin.Host == "" {
			origin.Host = oldestForwardedValue(fwd, "host")
		}
	}
	/* TODO: It doesn't make sense to support the X-Forwarded-X headers as long as the saml and oidc libraries don't support passing a custom issuer interceptor.
	if origin.Scheme == "" {
		origin.Scheme = r.Header.Get("X-Forwarded-Proto")
	}
	if origin.Host == "" {
		origin.Host = r.Header.Get("X-Forwarded-Host")
	}
	*/
	if origin.Scheme == "" {
		origin.Scheme = "http"
		if externalSecure {
			origin.Scheme = "https"
		}
	}
	if origin.Host == "" {
		origin.Host = r.Host
	}
	origin.Full = fmt.Sprintf("%s://%s", origin.Scheme, origin.Host)
	if !http_util.IsOrigin(origin.Full) {
		err = errors.Join(err, fmt.Errorf("invalid origin: %s", origin.Full))
	}
	var splitErr error
	origin.Domain, _, splitErr = net.SplitHostPort(origin.Host)
	if splitErr != nil && strings.Contains(splitErr.Error(), "missing port in address") {
		// We assume the schemes default port is not passed with the host, which is fine
		splitErr = nil
		origin.Domain = origin.Host
	}
	err = errors.Join(err, splitErr)
	return origin, err
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

// TODO: Do we have to check this?
// isAllowedToSendHTTP1Header check if the gRPC call was sent to `localhost`
// this is only possible when calling the server directly running on localhost
// or through the gRPC gateway
func isAllowedToSendHTTP1Header(md metadata.MD) bool {
	authority, ok := md[":authority"]
	return ok && len(authority) == 1 && strings.Split(authority[0], ":")[0] == "localhost"
}

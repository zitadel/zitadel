package tracing

import (
	"net/http"
	"strings"

	http_trace "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func shouldNotIgnore(endpoints ...string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		for _, endpoint := range endpoints {
			if strings.HasPrefix(r.URL.RequestURI(), endpoint) {
				return false
			}
		}
		return true
	}
}

func TraceHandler(handler http.Handler, ignoredEndpoints ...string) http.Handler {
	return http_trace.NewHandler(handler,
		"zitadel",
		http_trace.WithFilter(shouldNotIgnore(ignoredEndpoints...)),
		http_trace.WithPublicEndpoint(),
		http_trace.WithSpanNameFormatter(spanNameFormatter))
}

func spanNameFormatter(_ string, r *http.Request) string {
	return r.Host + r.URL.EscapedPath()
}

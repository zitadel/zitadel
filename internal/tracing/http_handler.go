package tracing

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func ignoredEndpointsFilter(endpoints ...string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		for _, endpoint := range endpoints {
			if strings.HasPrefix(r.URL.RequestURI(), endpoint) {
				return true
			}
		}
		return false
	}
}

func TraceHandler(handler http.Handler, ignoredEndpoints ...string) http.Handler {
	return otelhttp.NewHandler(handler, "zitadel", otelhttp.WithFilter(ignoredEndpointsFilter(ignoredEndpoints...)))
}

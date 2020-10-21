package tracing

import (
	"net/http"
	"strings"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func TraceHandler(handler http.Handler, ignoredEndpoints ...string) http.Handler {
	return &ochttp.Handler{
		Handler: handler,
		FormatSpanName: func(r *http.Request) string {
			host := r.URL.Host
			if host == "" {
				host = r.Host
			}
			return host + r.URL.Path
		},

		StartOptions: trace.StartOptions{Sampler: Sampler()},
		IsHealthEndpoint: func(r *http.Request) bool {
			for _, endpoint := range ignoredEndpoints {
				if strings.HasPrefix(r.URL.RequestURI(), endpoint) {
					return true
				}
			}
			return false
		},
	}
}

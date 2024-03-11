package telemetry

import (
	"net/http"
	"strings"

	"github.com/zitadel/zitadel/internal/telemetry/metrics"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

func TelemetryHandler(handler http.Handler, ignoredEndpoints ...string) http.Handler {
	return otelhttp.NewHandler(handler,
		"zitadel",
		otelhttp.WithFilter(shouldNotIgnore(ignoredEndpoints...)),
		otelhttp.WithPublicEndpoint(),
		otelhttp.WithMeterProvider(metrics.GetMetricsProvider()),
		otelhttp.WithServerName("ZITADEL"),
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
		otelhttp.WithPublicEndpoint(),
	)
}

func spanNameFormatter(_ string, r *http.Request) string {
	return r.URL.EscapedPath()
}

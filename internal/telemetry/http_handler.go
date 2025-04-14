package telemetry

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/zitadel/zitadel/internal/telemetry/metrics"
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
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
		otelhttp.WithMeterProvider(metrics.GetMetricsProvider()))
}

func spanNameFormatter(_ string, r *http.Request) string {
	return strings.Split(r.RequestURI, "?")[0]
}

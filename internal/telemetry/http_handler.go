package telemetry

import (
	"net/http"
	"slices"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

func shouldNotIgnore(endpoints ...string) func(r *http.Request) bool {
	return func(r *http.Request) bool {
		return !slices.ContainsFunc(endpoints, func(endpoint string) bool {
			return strings.HasPrefix(r.URL.RequestURI(), endpoint)
		})
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

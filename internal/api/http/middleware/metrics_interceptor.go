package middleware

import (
	"net/http"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
)

func DefaultMetricsHandler(handler http.Handler) http.Handler {
	metricTypes := []metrics.MetricType{metrics.MetricTypeTotalCount}
	return MetricsHandler(metricTypes, http_utils.Probes...)(handler)
}

func MetricsHandler(metricTypes []metrics.MetricType, ignoredMethods ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return metrics.NewMetricsHandler(handler, metricTypes, ignoredMethods...)
	}
}

package middleware

import (
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"net/http"

	http_utils "github.com/caos/zitadel/internal/api/http"
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

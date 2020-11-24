package middleware

import (
	"github.com/caos/zitadel/internal/telemetry/metrics"
	"net/http"

	http_utils "github.com/caos/zitadel/internal/api/http"
)

func DefaultMetricsyHandler(handler http.Handler) http.Handler {
	return MetricsHandler(http_utils.Probes...)(handler)
}

func MetricsHandler(ignoredMethods ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return metrics.NewMetricsHandler(handler, ignoredMethods...)
	}
}

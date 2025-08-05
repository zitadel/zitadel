package middleware

import (
	"net/http"

	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/telemetry"
)

func DefaultTelemetryHandler(handler http.Handler) http.Handler {
	return TelemetryHandler(http_utils.Probes...)(handler)
}

func TelemetryHandler(ignoredMethods ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return telemetry.TelemetryHandler(handler, ignoredMethods...)
	}
}

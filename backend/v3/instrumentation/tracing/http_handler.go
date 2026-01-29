package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
)

func NewHandler(handler http.Handler, ignoredPrefix ...string) http.Handler {
	return otelhttp.NewHandler(handler,
		"zitadel",
		otelhttp.WithFilter(instrumentation.RequestFilter(ignoredPrefix...)),
		otelhttp.WithPublicEndpoint(),
		otelhttp.WithSpanNameFormatter(spanNameFormatter),
		otelhttp.WithMeterProvider(otel.GetMeterProvider()))
}

func spanNameFormatter(_ string, r *http.Request) string {
	return r.URL.Path
}

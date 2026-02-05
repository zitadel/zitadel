package middleware

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

func DefaultTraceHandler(handler http.Handler) http.Handler {
	return TraceHandler(http_utils.Probes...)(handler)
}

func TraceHandler(ignoredPrefix ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return otelhttp.NewHandler(handler,
			"zitadel",
			otelhttp.WithFilter(instrumentation.RequestFilter(ignoredPrefix...)),
			otelhttp.WithPublicEndpoint(),
			otelhttp.WithSpanNameFormatter(spanNameFormatter),
			otelhttp.WithMeterProvider(otel.GetMeterProvider()),
		)
	}
}

func spanNameFormatter(_ string, r *http.Request) string {
	return r.URL.Path
}

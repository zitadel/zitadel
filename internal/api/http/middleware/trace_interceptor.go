package middleware

import (
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/tracing"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
)

func DefaultTraceHandler(handler http.Handler) http.Handler {
	return TraceHandler(http_utils.Probes...)(handler)
}

func TraceHandler(ignoredPrefix ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return tracing.NewHandler(handler, ignoredPrefix...)
	}
}

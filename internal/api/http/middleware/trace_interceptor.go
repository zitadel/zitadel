package middleware

import (
	"net/http"

	"github.com/caos/zitadel/internal/api"
	"github.com/caos/zitadel/internal/tracing"
)

func DefaultTraceHandler(handler http.Handler) http.Handler {
	return tracing.TraceHandler(handler, api.Probes...)
}

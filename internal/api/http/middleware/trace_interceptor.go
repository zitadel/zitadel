package middleware

import (
	"net/http"

	http2 "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/tracing"
)

func DefaultTraceHandler(handler http.Handler) http.Handler {
	return tracing.TraceHandler(handler, http2.Probes...)
}

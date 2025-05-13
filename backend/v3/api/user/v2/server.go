package userv2

import (
	"github.com/zitadel/zitadel/backend/v3/telemetry/logging"
	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
)

var (
	logger logging.Logger
	tracer tracing.Tracer
)

func SetLogger(l logging.Logger) {
	logger = l
}

func SetTracer(t tracing.Tracer) {
	tracer = t
}

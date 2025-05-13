package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Tracer struct {
	trace.Tracer
}

var noopTracer = Tracer{
	Tracer: noop.NewTracerProvider().Tracer(""),
}

func (t *Tracer) Start(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if t.Tracer == nil {
		return noopTracer.Start(ctx, spanName, opts...)
	}
	return t.Tracer.Start(ctx, spanName, opts...)
}

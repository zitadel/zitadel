package otel

import (
	"context"
	"net/http"

	"github.com/caos/zitadel/internal/telemetry/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	export "go.opentelemetry.io/otel/sdk/export/trace"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Tracer struct {
	Exporter trace.Tracer
	sampler  sdk_trace.Sampler
}

func NewTracer(name string, sampler sdk_trace.Sampler, exporter export.SpanExporter) *Tracer {
	tp := sdk_trace.NewTracerProvider(
		sdk_trace.WithSampler(sampler),
		sdk_trace.WithSyncer(exporter),
	)

	otel.SetTracerProvider(tp)
	tc := propagation.TraceContext{}
	otel.SetTextMapPropagator(tc)

	return &Tracer{Exporter: tp.Tracer(name), sampler: sampler}
}

func (t *Tracer) Sampler() sdk_trace.Sampler {
	return t.sampler
}

func (t *Tracer) NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
}

func (t *Tracer) NewServerSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, trace.WithSpanKind(trace.SpanKindServer))
}

func (t *Tracer) NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, name, trace.WithSpanKind(trace.SpanKindClient))
}

func (t *Tracer) NewClientSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, trace.WithSpanKind(trace.SpanKindClient))
}

func (t *Tracer) NewSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller)
}

func (t *Tracer) newSpan(ctx context.Context, caller string, options ...trace.SpanOption) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, caller, options...)
}

func (t *Tracer) newSpanFromName(ctx context.Context, name string, options ...trace.SpanOption) (context.Context, *tracing.Span) {
	ctx, span := t.Exporter.Start(ctx, name, options...)
	return ctx, tracing.CreateSpan(span)
}

func (t *Tracer) NewSpanHTTP(r *http.Request, caller string) (*http.Request, *tracing.Span) {
	ctx, span := t.NewSpan(r.Context(), caller)
	r = r.WithContext(ctx)
	return r, span
}

package otel

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	api_trace "go.opentelemetry.io/otel/trace"

	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type Tracer struct {
	Exporter api_trace.Tracer
	sampler  sdk_trace.Sampler
}

func NewTracer(name string, sampler sdk_trace.Sampler, exporter sdk_trace.SpanExporter) *Tracer {
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
	return t.newSpanFromName(ctx, name, api_trace.WithSpanKind(api_trace.SpanKindServer))
}

func (t *Tracer) NewServerSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, api_trace.WithSpanKind(api_trace.SpanKindServer))
}

func (t *Tracer) NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, name, api_trace.WithSpanKind(api_trace.SpanKindClient))
}

func (t *Tracer) NewClientSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, api_trace.WithSpanKind(api_trace.SpanKindClient))
}

func (t *Tracer) NewSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller)
}

func (t *Tracer) newSpan(ctx context.Context, caller string, options ...api_trace.SpanStartOption) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, caller, options...)
}

func (t *Tracer) newSpanFromName(ctx context.Context, name string, options ...api_trace.SpanStartOption) (context.Context, *tracing.Span) {
	ctx, span := t.Exporter.Start(ctx, name, options...)
	return ctx, tracing.CreateSpan(span)
}

func (t *Tracer) NewSpanHTTP(r *http.Request, caller string) (*http.Request, *tracing.Span) {
	ctx, span := t.NewSpan(r.Context(), caller)
	r = r.WithContext(ctx)
	return r, span
}

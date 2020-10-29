package otel

import (
	"context"
	"net/http"

	"github.com/caos/zitadel/internal/tracing"
	"go.opentelemetry.io/otel/api/global"
	apitrace "go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/sdk/export/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Tracer struct {
	Exporter apitrace.Tracer
	sampler  sdktrace.Sampler
}

func NewTracer(name string, sampler sdktrace.Sampler, exporter trace.SpanExporter) *Tracer {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sampler}),
		sdktrace.WithSyncer(exporter),
	)

	global.SetTracerProvider(tp)

	return &Tracer{Exporter: tp.Tracer(name), sampler: sampler}
}

func (t *Tracer) Sampler() sdktrace.Sampler {
	return t.sampler
}

func (t *Tracer) NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, name, apitrace.WithSpanKind(apitrace.SpanKindServer))
}

func (t *Tracer) NewServerSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, apitrace.WithSpanKind(apitrace.SpanKindServer))
}

func (t *Tracer) NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, name, apitrace.WithSpanKind(apitrace.SpanKindClient))
}

func (t *Tracer) NewClientSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller, apitrace.WithSpanKind(apitrace.SpanKindClient))
}

func (t *Tracer) NewSpan(ctx context.Context, caller string) (context.Context, *tracing.Span) {
	return t.newSpan(ctx, caller)
}

func (t *Tracer) newSpan(ctx context.Context, caller string, options ...apitrace.SpanOption) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, caller, options...)
}

func (t *Tracer) newSpanFromName(ctx context.Context, name string, options ...apitrace.SpanOption) (context.Context, *tracing.Span) {
	//TODO: should we use apitrace.SpanFromContext?

	ctx, span := t.Exporter.Start(ctx, name, options...)
	return ctx, tracing.CreateSpan(span)
}

func (t *Tracer) NewSpanHTTP(r *http.Request, caller string) (*http.Request, *tracing.Span) {
	ctx, span := t.NewSpan(r.Context(), caller)
	r = r.WithContext(ctx)
	return r, span
}

func (t *Tracer) SetErrStatus(span apitrace.Span, code int32, err error, obj ...string) {
	span.RecordError(context.TODO(), err, apitrace.WithErrorStatus(codes.Error))
	span.SetAttributes(label.Int32("code", code), label.Array("error_attributes", obj))
}

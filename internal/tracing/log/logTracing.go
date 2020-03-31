package log

import (
	"context"
	"net/http"

	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/tracing"
)

type Tracer struct {
	sampler trace.Sampler
}

func (t *Tracer) Start() error {
	trace.RegisterExporter(&exporter.PrintExporter{})

	views := append(ocgrpc.DefaultServerViews, ocgrpc.DefaultClientViews...)
	views = append(views, ochttp.DefaultClientViews...)
	views = append(views, ochttp.DefaultServerViews...)

	if err := view.Register(views...); err != nil {
		return errors.ThrowInternal(err, "LOG-PoFiB", "unable to register view")
	}

	trace.ApplyConfig(trace.Config{DefaultSampler: t.sampler})

	return nil
}

func (t *Tracer) Sampler() trace.Sampler {
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

func (t *Tracer) newSpan(ctx context.Context, caller string, options ...trace.StartOption) (context.Context, *tracing.Span) {
	return t.newSpanFromName(ctx, caller, options...)
}

func (t *Tracer) newSpanFromName(ctx context.Context, name string, options ...trace.StartOption) (context.Context, *tracing.Span) {
	ctx, span := trace.StartSpan(ctx, name, options...)
	return ctx, tracing.CreateSpan(span)
}

func (t *Tracer) NewSpanHTTP(r *http.Request, caller string) (*http.Request, *tracing.Span) {
	ctx, span := t.NewSpan(r.Context(), caller)
	r = r.WithContext(ctx)
	return r, span
}

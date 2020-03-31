package google

import (
	"context"
	"net/http"
	"os"
	"strings"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/tracing"
)

type Tracer struct {
	Exporter     *stackdriver.Exporter
	projectID    string
	metricPrefix string
	sampler      trace.Sampler
}

func (t *Tracer) Start() (err error) {
	t.Exporter, err = stackdriver.NewExporter(stackdriver.Options{
		ProjectID:    t.projectID,
		MetricPrefix: t.metricPrefix,
	})
	if err != nil {
		return errors.ThrowInternal(err, "GOOGL-4dCnX", "unable to start exporter")
	}

	views := append(ocgrpc.DefaultServerViews, ocgrpc.DefaultClientViews...)
	views = append(views, ochttp.DefaultClientViews...)
	views = append(views, ochttp.DefaultServerViews...)

	if err = view.Register(views...); err != nil {
		return errors.ThrowInternal(err, "GOOGL-Q6L6w", "unable to register view")
	}

	trace.RegisterExporter(t.Exporter)
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

func envIsSet() bool {
	gAuthCred := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	return strings.Contains(gAuthCred, ".json")
}

func (t *Tracer) SetErrStatus(span *trace.Span, code int32, err error, obj ...string) {
	span.SetStatus(trace.Status{Code: code, Message: err.Error() + strings.Join(obj, ", ")})
}

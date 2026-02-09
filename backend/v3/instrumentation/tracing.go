package instrumentation

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"

	google_trace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type TraceConfig struct {
	Fraction         float64
	Exporter         ExporterConfig
	TrustRemoteSpans bool
}

// TODO: remove for v5 release
type LegacyTraceConfig struct {
	Type      string
	Fraction  float64
	Endpoint  string
	ProjectID string
}

// SetLegacyConfig sets the fields of the TraceConfig based on the provided LegacyTraceConfig.
// If the Type field is already set to a value other than [ExporterTypeNone], the legacy config is ignored.
func (c *TraceConfig) SetLegacyConfig(lc *LegacyTraceConfig) {
	typ := c.Exporter.Type
	if lc == nil || !typ.isNone() {
		return
	}
	typeMap := map[string]ExporterType{
		"otel":   ExporterTypeGRPC,
		"google": ExporterTypeGoogle,
		"log":    ExporterTypeStdOut,
	}

	c.Fraction = lc.Fraction
	c.Exporter.Type = typeMap[lc.Type]
	c.Exporter.Endpoint = lc.Endpoint
	c.Exporter.GoogleProjectID = lc.ProjectID
}

type Tracer struct {
	trace.Tracer
}

func (t *Tracer) NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *Span) {
	return t.newSpanFromName(ctx, name, trace.WithSpanKind(trace.SpanKindServer))
}

func (t *Tracer) NewServerSpan(ctx context.Context, caller string) (context.Context, *Span) {
	return t.newSpan(ctx, caller, trace.WithSpanKind(trace.SpanKindServer))
}

func (t *Tracer) NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *Span) {
	return t.newSpanFromName(ctx, name, trace.WithSpanKind(trace.SpanKindClient))
}

func (t *Tracer) NewClientSpan(ctx context.Context, caller string) (context.Context, *Span) {
	return t.newSpan(ctx, caller, trace.WithSpanKind(trace.SpanKindClient))
}

func (t *Tracer) NewSpan(ctx context.Context, caller string) (context.Context, *Span) {
	return t.newSpan(ctx, caller)
}

func (t *Tracer) newSpan(ctx context.Context, caller string, options ...trace.SpanStartOption) (context.Context, *Span) {
	return t.newSpanFromName(ctx, caller, options...)
}

func (t *Tracer) newSpanFromName(ctx context.Context, name string, options ...trace.SpanStartOption) (context.Context, *Span) {
	ctx, span := t.Start(ctx, name, options...)
	return ctx, CreateSpan(span)
}

func (t *Tracer) NewSpanHTTP(r *http.Request, caller string) (*http.Request, *Span) {
	ctx, span := t.NewSpan(r.Context(), caller)
	r = r.WithContext(ctx)
	return r, span
}

type Span struct {
	span trace.Span
	opts []trace.SpanEndOption
}

func CreateSpan(span trace.Span) *Span {
	return &Span{span: span, opts: []trace.SpanEndOption{}}
}

func (s *Span) End() {
	if s.span == nil {
		return
	}

	s.span.End(s.opts...)
}

func (s *Span) EndWithError(err error) {
	s.SetStatusByError(err)
	s.End()
}

func (s *Span) SetStatusByError(err error) {
	if s.span == nil {
		return
	}
	if err != nil {
		s.span.RecordError(err)
		s.span.SetStatus(codes.Error, err.Error())
	}
	var zerr *zerrors.ZitadelError
	if errors.As(err, &zerr) {
		s.span.SetAttributes(
			attribute.Stringer("error_kind", zerr.Kind),
			attribute.String("error_msg", zerr.Message),
			attribute.String("error_id", zerr.ID),
		)
	}
}

func newTracerProvider(ctx context.Context, cfg TraceConfig, resource *resource.Resource) (_ *sdk_trace.TracerProvider, err error) {
	var exporter sdk_trace.SpanExporter
	switch cfg.Exporter.Type {
	case ExporterTypeUnspecified, ExporterTypeNone:
		// no exporter
	case ExporterTypeStdOut, ExporterTypeStdErr:
		exporter, err = traceStdOutExporter(cfg.Exporter)
	case ExporterTypeGRPC:
		exporter, err = traceGrpcExporter(ctx, cfg.Exporter)
	case ExporterTypeHTTP:
		exporter, err = traceHttpExporter(ctx, cfg.Exporter)
	case ExporterTypeGoogle:
		exporter, err = traceGoogleExporter(ctx, cfg.Exporter)
	case ExporterTypePrometheus:
		fallthrough // prometheus is not supported for logs
	default:
		err = errExporterType(cfg.Exporter.Type, "tracer")
	}
	if err != nil {
		return nil, fmt.Errorf("trace exporter: %w", err)
	}
	sampler := newSampler(cfg)
	opts := []sdk_trace.TracerProviderOption{
		sdk_trace.WithResource(resource),
		sdk_trace.WithSampler(sampler),
	}
	if exporter != nil {
		opts = append(opts, sdk_trace.WithBatcher(exporter))
	}
	tracerProvider := sdk_trace.NewTracerProvider(opts...)
	return tracerProvider, nil
}

func traceStdOutExporter(cfg ExporterConfig) (sdk_trace.SpanExporter, error) {
	options := []stdouttrace.Option{
		stdouttrace.WithPrettyPrint(),
	}
	if cfg.Type == ExporterTypeStdErr {
		options = append(options, stdouttrace.WithWriter(os.Stderr))
	}
	exporter, err := stdouttrace.New(options...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func traceGrpcExporter(ctx context.Context, cfg ExporterConfig) (sdk_trace.SpanExporter, error) {
	var grpcOpts []otlptracegrpc.Option
	if cfg.Endpoint != "" {
		grpcOpts = append(grpcOpts, otlptracegrpc.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		grpcOpts = append(grpcOpts, otlptracegrpc.WithInsecure())
	}

	exporter, err := otlptracegrpc.New(ctx, grpcOpts...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func traceHttpExporter(ctx context.Context, cfg ExporterConfig) (sdk_trace.SpanExporter, error) {
	var httpOpts []otlptracehttp.Option
	if cfg.Endpoint != "" {
		httpOpts = append(httpOpts, otlptracehttp.WithEndpoint(cfg.Endpoint))
	}
	if cfg.Insecure {
		httpOpts = append(httpOpts, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, httpOpts...)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

func traceGoogleExporter(ctx context.Context, cfg ExporterConfig) (sdk_trace.SpanExporter, error) {
	exporter, err := google_trace.New(
		google_trace.WithContext(ctx),
		google_trace.WithProjectID(cfg.GoogleProjectID),
	)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}

// newSampler returns a sampler decorator which behaves differently,
// based on the parent of the span. If the span has no parent and is of kind server,
// the decorated sampler is used to make sampling decision.
// If the span has a parent, depending on whether the parent is remote and whether it
// is sampled, one of the following samplers will apply:
//   - remote parent sampled -> always sample
//   - remote parent not sampled -> sample based on the decorated sampler (fraction based)
//   - local parent sampled -> always sample
//   - local parent not sampled -> never sample
func newSampler(cfg TraceConfig) sdk_trace.Sampler {
	fraction := sdk_trace.TraceIDRatioBased(cfg.Fraction)
	return sdk_trace.ParentBased(
		spanKindBased(fraction, trace.SpanKindServer),
		sdk_trace.WithRemoteParentNotSampled(fraction),
	)
}

type spanKindSampler struct {
	sampler sdk_trace.Sampler
	kinds   []trace.SpanKind
}

// ShouldSample implements the [sdk_trace.Sampler] interface.
// It will not sample any spans which do not match the configured span kinds.
// For spans which do match, the decorated sampler is used to make the sampling decision.
func (sk spanKindSampler) ShouldSample(p sdk_trace.SamplingParameters) sdk_trace.SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	if !slices.Contains(sk.kinds, p.Kind) {
		return sdk_trace.SamplingResult{
			Decision:   sdk_trace.Drop,
			Tracestate: psc.TraceState(),
		}
	}
	s := sk.sampler.ShouldSample(p)
	return s
}

func (sk spanKindSampler) Description() string {
	return fmt.Sprintf("SpanKindBased{sampler:%s,kinds:%v}",
		sk.sampler.Description(),
		sk.kinds,
	)
}

// spanKindBased returns a sampler decorator which behaves differently, based on the kind of the span.
// If the span kind does not match one of the configured kinds, it will not be sampled.
// If the span kind matches, the decorated sampler is used to make sampling decision.
func spanKindBased(sampler sdk_trace.Sampler, kinds ...trace.SpanKind) sdk_trace.Sampler {
	return spanKindSampler{
		sampler: sampler,
		kinds:   kinds,
	}
}

package tracing

import (
	"context"
	"net/http"

	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Tracer interface {
	NewSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewClientSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewServerSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *Span)
	NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *Span)
	NewSpanHTTP(r *http.Request, caller string) (*http.Request, *Span)
	Sampler() sdk_trace.Sampler
}

type Config interface {
	NewTracer() error
}

var T Tracer

func Sampler() sdk_trace.Sampler {
	if T == nil {
		return sdk_trace.NeverSample()
	}
	return T.Sampler()
}

func NewSpan(ctx context.Context) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewSpan(ctx, GetCaller())
}

func NewNamedSpan(ctx context.Context, name string) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewSpan(ctx, name)
}

func NewClientSpan(ctx context.Context) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewClientSpan(ctx, GetCaller())
}

func NewServerSpan(ctx context.Context) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewServerSpan(ctx, GetCaller())
}

func NewClientInterceptorSpan(ctx context.Context) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewClientInterceptorSpan(ctx, GetCaller())
}

func NewServerInterceptorSpan(ctx context.Context) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewServerInterceptorSpan(ctx, GetCaller())
}

func NewSpanHTTP(r *http.Request) (*http.Request, *Span) {
	if T == nil {
		return r, CreateSpan(nil)
	}
	return T.NewSpanHTTP(r, GetCaller())
}

func TraceIDFromCtx(ctx context.Context) string {
	return trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}

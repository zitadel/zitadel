package tracing

import (
	"context"
	"net/http"

	"go.opencensus.io/trace"
)

type Tracer interface {
	Start() error
	NewSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewClientSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewServerSpan(ctx context.Context, caller string) (context.Context, *Span)
	NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *Span)
	NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *Span)
	NewSpanHTTP(r *http.Request, caller string) (*http.Request, *Span)
	Sampler() trace.Sampler
}

type Config interface {
	NewTracer() error
}

var T Tracer

func Sampler() trace.Sampler {
	if T == nil {
		return trace.NeverSample()
	}
	return T.Sampler()
}

func NewSpan(ctx context.Context) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewSpan(ctx, GetCaller())
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

func NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewClientInterceptorSpan(ctx, name)
}

func NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *Span) {
	if T == nil {
		return ctx, CreateSpan(nil)
	}
	return T.NewServerInterceptorSpan(ctx, name)
}

func NewSpanHTTP(r *http.Request) (*http.Request, *Span) {
	if T == nil {
		return r, CreateSpan(nil)
	}
	return T.NewSpanHTTP(r, GetCaller())
}

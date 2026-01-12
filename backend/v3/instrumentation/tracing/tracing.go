// Package tracing provides helper functions to create spans for telemetry tracing.
package tracing

import (
	"context"
	"net/http"
	"sync"

	api_trace "go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/cmd/build"
)

type Tracer interface {
	NewSpan(ctx context.Context, caller string) (context.Context, *instrumentation.Span)
	NewClientSpan(ctx context.Context, caller string) (context.Context, *instrumentation.Span)
	NewServerSpan(ctx context.Context, caller string) (context.Context, *instrumentation.Span)
	NewClientInterceptorSpan(ctx context.Context, name string) (context.Context, *instrumentation.Span)
	NewServerInterceptorSpan(ctx context.Context, name string) (context.Context, *instrumentation.Span)
	NewSpanHTTP(r *http.Request, caller string) (*http.Request, *instrumentation.Span)
}

var globalTracer = sync.OnceValue(func() Tracer {
	return instrumentation.NewTracer(
		instrumentation.Name,
		api_trace.WithInstrumentationVersion(build.Version()),
	)
})

func GlobalTracer() Tracer {
	return globalTracer()
}

func NewSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewNamedSpan(ctx context.Context, name string) (context.Context, *instrumentation.Span) {
	return globalTracer().NewSpan(ctx, name)
}

func NewClientSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewClientSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewServerSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewServerSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewClientInterceptorSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewClientInterceptorSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewServerInterceptorSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewServerInterceptorSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewSpanHTTP(r *http.Request) (*http.Request, *instrumentation.Span) {
	return globalTracer().NewSpanHTTP(r, instrumentation.GetCallingFunc(1))
}

func TraceIDFromCtx(ctx context.Context) string {
	return api_trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}

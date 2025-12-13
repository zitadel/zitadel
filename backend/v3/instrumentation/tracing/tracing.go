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

const pkgName = "github.com/zitadel/zitadel/backend/v3/instrumentation/tracing"

var globalTracer = sync.OnceValue(func() Tracer {
	return instrumentation.NewTracer(
		pkgName,
		api_trace.WithInstrumentationVersion(build.Version()),
	)
})

func NewSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewSpan(ctx, GetCaller())
}

func NewNamedSpan(ctx context.Context, name string) (context.Context, *instrumentation.Span) {
	return globalTracer().NewSpan(ctx, name)
}

func NewClientSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewClientSpan(ctx, GetCaller())
}

func NewServerSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewServerSpan(ctx, GetCaller())
}

func NewClientInterceptorSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewClientInterceptorSpan(ctx, GetCaller())
}

func NewServerInterceptorSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return globalTracer().NewServerInterceptorSpan(ctx, GetCaller())
}

func NewSpanHTTP(r *http.Request) (*http.Request, *instrumentation.Span) {
	return globalTracer().NewSpanHTTP(r, GetCaller())
}

func TraceIDFromCtx(ctx context.Context) string {
	return api_trace.SpanFromContext(ctx).SpanContext().TraceID().String()
}

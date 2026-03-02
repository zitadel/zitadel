// Package tracing acts as a shim to [tracing],
// so we can migrate the tracing package without changing all imports at once.
package tracing

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/backend/v3/instrumentation"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/tracing"
)

func NewSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return tracing.GlobalTracer().NewSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewNamedSpan(ctx context.Context, name string) (context.Context, *instrumentation.Span) {
	return tracing.NewNamedSpan(ctx, name)
}

func NewClientSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return tracing.GlobalTracer().NewClientSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewServerSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return tracing.GlobalTracer().NewServerSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewClientInterceptorSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return tracing.GlobalTracer().NewClientInterceptorSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewServerInterceptorSpan(ctx context.Context) (context.Context, *instrumentation.Span) {
	return tracing.GlobalTracer().NewServerInterceptorSpan(ctx, instrumentation.GetCallingFunc(1))
}

func NewSpanHTTP(r *http.Request) (*http.Request, *instrumentation.Span) {
	return tracing.GlobalTracer().NewSpanHTTP(r, instrumentation.GetCallingFunc(1))
}

func TraceIDFromCtx(ctx context.Context) string {
	return tracing.TraceIDFromCtx(ctx)
}

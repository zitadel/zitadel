package info

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func HTTPPathIntoContext(path string) runtime.AnnotateContextOption {
	return runtime.WithHTTPPathPattern(path)
}

func HTTPPathFromContext() func(context.Context) (string, bool) {
	return runtime.HTTPPathPattern
}

func RPCMethodFromContext() func(context.Context) (string, bool) {
	return runtime.RPCMethod
}

type requestMethodKey struct{}

func RequestMethodIntoContext(method string) runtime.AnnotateContextOption {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, requestMethodKey{}, method)
	}
}

func RequestMethodFromContext() func(context.Context) (string, bool) {
	return func(ctx context.Context) (string, bool) {
		m := ctx.Value(requestMethodKey{})
		if m == nil {
			return "", false
		}
		ms, ok := m.(string)
		if !ok {
			return "", false
		}
		return ms, true
	}
}

package handler

import (
	"context"
)

// Handler is a function that handles the request.
type Handler[Req, Res any] func(ctx context.Context, request Req) (res Res, err error)

// Decorator is a function that decorates the handle function.
type Decorator[Req, Res any] func(ctx context.Context, request Req, handle Handler[Req, Res]) (res Res, err error)

// Chain chains the handle function with the next handler.
// The next handler is called after the handle function.
func Chain[Req, Res any](handle Handler[Req, Res], next Handler[Res, Res]) Handler[Req, Res] {
	return func(ctx context.Context, request Req) (res Res, err error) {
		res, err = handle(ctx, request)
		if err != nil {
			return res, err
		}
		return next(ctx, res)
	}
}

func Chains[Req, Res any](handle Handler[Req, Res], nexts ...Handler[Res, Res]) Handler[Req, Res] {
	return func(ctx context.Context, request Req) (res Res, err error) {
		for _, next := range nexts {
			handle = Chain(handle, next)
		}
		return handle(ctx, request)
	}
}

// Decorate decorates the handle function with the decorate function.
// The decorate function is called before the handle function.
func Decorate[Req, Res any](handle Handler[Req, Res], decorate Decorator[Req, Res]) Handler[Req, Res] {
	return func(ctx context.Context, request Req) (res Res, err error) {
		return decorate(ctx, request, handle)
	}
}

// Decorates decorates the handle function with the decorate functions.
// The decorates function is called before the handle function.
func Decorates[Req, Res any](handle Handler[Req, Res], decorates ...Decorator[Req, Res]) Handler[Req, Res] {
	return func(ctx context.Context, request Req) (res Res, err error) {
		for i := len(decorates) - 1; i >= 0; i-- {
			handle = Decorate(handle, decorates[i])
		}
		return handle(ctx, request)
	}
}

// SkipNext skips the next handler if the handle function returns a non-nil response.
func SkipNext[Req, Res any](handle Handler[Req, Res], next Handler[Req, Res]) Handler[Req, Res] {
	return func(ctx context.Context, request Req) (res Res, err error) {
		var empty Res
		res, err = handle(ctx, request)
		// TODO: does this work?
		if any(res) == any(empty) || err != nil {
			return res, err
		}
		return next(ctx, request)
	}
}

// SkipNilHandler skips the handle function if the handler is nil.
// The function is safe to call with nil handler.
func SkipNilHandler[O, R any](handler *O, handle Handler[R, R]) Handler[R, R] {
	return func(ctx context.Context, request R) (res R, err error) {
		if handler == nil {
			return request, nil
		}
		return handle(ctx, request)
	}
}

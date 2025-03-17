package handler

import (
	"context"
)

// Handler is a function that handles the in.
type Handler[Out, In any] func(ctx context.Context, in Out) (out In, err error)

// Decorator is a function that decorates the handle function.
type Decorator[In, Out any] func(ctx context.Context, in In, handle Handler[In, Out]) (out Out, err error)

// Chain chains the handle function with the next handler.
// The next handler is called after the handle function.
func Chain[In, Out any](handle Handler[In, Out], next Handler[Out, Out]) Handler[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		out, err = handle(ctx, in)
		if err != nil {
			return out, err
		}
		return next(ctx, out)
	}
}

func Chains[In, Out any](handle Handler[In, Out], nexts ...Handler[Out, Out]) Handler[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		for _, next := range nexts {
			handle = Chain(handle, next)
		}
		return handle(ctx, in)
	}
}

// Decorate decorates the handle function with the decorate function.
// The decorate function is called before the handle function.
func Decorate[In, Out any](handle Handler[In, Out], decorate Decorator[In, Out]) Handler[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		return decorate(ctx, in, handle)
	}
}

// Decorates decorates the handle function with the decorate functions.
// The decorates function is called before the handle function.
func Decorates[In, Out any](handle Handler[In, Out], decorates ...Decorator[In, Out]) Handler[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		for i := len(decorates) - 1; i >= 0; i-- {
			handle = Decorate(handle, decorates[i])
		}
		return handle(ctx, in)
	}
}

// SkipNext skips the next handler if the handle function returns a non-nil response.
func SkipNext[In, Out any](handle Handler[In, Out], next Handler[In, Out]) Handler[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		var empty Out
		out, err = handle(ctx, in)
		// TODO: does this work?
		if any(out) != any(empty) || err != nil {
			return out, err
		}
		return next(ctx, in)
	}
}

// SkipNilHandler skips the handle function if the handler is nil.
// If handle is nil, an empty output is returned.
// The function is safe to call with nil handler.
func SkipNilHandler[O, In, Out any](handler *O, handle Handler[In, Out]) Handler[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		if handler == nil {
			return out, nil
		}
		return handle(ctx, in)
	}
}

// SkipReturnPreviousHandler skips the handle function if the handler is nil and returns the input.
// The function is safe to call with nil handler.
func SkipReturnPreviousHandler[O, In any](handler *O, handle Handler[In, In]) Handler[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		if handler == nil {
			return in, nil
		}
		return handle(ctx, in)
	}
}

func ResFuncToHandle[In any, Out any](fn func(context.Context, In) Out) Handler[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		return fn(ctx, in), nil
	}
}

func ErrFuncToHandle[In any](fn func(context.Context, In) error) Handler[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		err = fn(ctx, in)
		if err != nil {
			return out, err
		}
		return in, nil
	}
}

func NoReturnToHandle[In any](fn func(context.Context, In)) Handler[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		fn(ctx, in)
		return in, nil
	}
}

package handler

import (
	"context"
)

type Parameter[P, C any] struct {
	Previous P
	Current  C
}

// Handle is a function that handles the in.
type Handle[In, Out any] func(ctx context.Context, in In) (out Out, err error)

type DeferrableHandle[In, Out any] func(ctx context.Context, in In) (out Out, deferrable func(context.Context, error) error, err error)

type Defer[In, Out, NextOut any] func(handle DeferrableHandle[In, Out], next Handle[Out, NextOut]) Handle[In, NextOut]

type HandleNoReturn[In any] func(ctx context.Context, in In) error

// Middleware is a function that decorates the handle function.
// It must call the handle function but its up the the middleware to decide when and how.
type Middleware[In, Out any] func(ctx context.Context, in In, handle Handle[In, Out]) (out Out, err error)

func Deferrable[In, Out, NextOut any](handle DeferrableHandle[In, Out], next Handle[Out, NextOut]) Handle[In, NextOut] {
	return func(ctx context.Context, in In) (nextOut NextOut, err error) {
		out, deferrable, err := handle(ctx, in)
		if err != nil {
			return nextOut, err
		}
		defer func() {
			err = deferrable(ctx, err)
		}()
		return next(ctx, out)
	}
}

// Chain chains the handle function with the next handler.
// The next handler is called after the handle function.
func Chain[In, Out, NextOut any](handle Handle[In, Out], next Handle[Out, NextOut]) Handle[In, NextOut] {
	return func(ctx context.Context, in In) (nextOut NextOut, err error) {
		out, err := handle(ctx, in)
		if err != nil {
			return nextOut, err
		}
		return next(ctx, out)
	}
}

// Chains chains the handle function with the next handlers.
// The next handlers are called after the handle function.
// The order of the handlers is preserved.
func Chains[In, Out any](handle Handle[In, Out], chain ...Handle[Out, Out]) Handle[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		for _, next := range chain {
			handle = Chain(handle, next)
		}
		return handle(ctx, in)
	}
}

// Decorate decorates the handle function with the decorate function.
// The decorate function is called before the handle function.
func Decorate[In, Out any](handle Handle[In, Out], decorate Middleware[In, Out]) Handle[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		return decorate(ctx, in, handle)
	}
}

// Decorates decorates the handle function with the decorate functions.
// The decorates function is called before the handle function.
func Decorates[In, Out any](handle Handle[In, Out], decorates ...Middleware[In, Out]) Handle[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		for i := len(decorates) - 1; i >= 0; i-- {
			handle = Decorate(handle, decorates[i])
		}
		return handle(ctx, in)
	}
}

// SkipNext skips the next handler if the handle function returns a non-empty output or an error.
func SkipNext[In, Out any](handle Handle[In, Out], next Handle[In, Out]) Handle[In, Out] {
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

func HandleIf[In any](cond func(In) bool, handle Handle[In, In]) Handle[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		if !cond(in) {
			return in, nil
		}
		return handle(ctx, in)
	}
}

func SkipIf[In any](cond func(In) bool, handle Handle[In, In]) Handle[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		if cond(in) {
			return in, nil
		}
		return handle(ctx, in)
	}
}

// SkipNilHandler skips the handle function if the handler is nil.
// If handle is nil, an empty output is returned.
// The function is safe to call with nil handler.
func SkipNilHandler[O, In, Out any](handler *O, handle Handle[In, Out]) Handle[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		if handler == nil {
			return out, nil
		}
		return handle(ctx, in)
	}
}

// SkipReturnPreviousHandler skips the handle function if the handler is nil and returns the input.
// The function is safe to call with nil handler.
func SkipReturnPreviousHandler[O, In any](handler *O, handle Handle[In, In]) Handle[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		if handler == nil {
			return in, nil
		}
		return handle(ctx, in)
	}
}

func CtxFuncToHandle[Out any](fn func(context.Context) (Out, error)) Handle[struct{}, Out] {
	return func(ctx context.Context, in struct{}) (out Out, err error) {
		return fn(ctx)
	}
}

func ResFuncToHandle[In any, Out any](fn func(context.Context, In) Out) Handle[In, Out] {
	return func(ctx context.Context, in In) (out Out, err error) {
		return fn(ctx, in), nil
	}
}

func ErrFuncToHandle[In any](fn func(context.Context, In) error) Handle[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		err = fn(ctx, in)
		if err != nil {
			return out, err
		}
		return in, nil
	}
}

func NoReturnToHandle[In any](fn func(context.Context, In)) Handle[In, In] {
	return func(ctx context.Context, in In) (out In, err error) {
		fn(ctx, in)
		return in, nil
	}
}

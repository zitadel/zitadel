package domain

import (
	"context"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type lazyGetter[T any] struct {
	get        func(ctx context.Context, opts *InvokeOpts) (T, error)
	wasFetched bool
	value      T
	err        error
}

func (f *lazyGetter[T]) fetch(ctx context.Context, opts *InvokeOpts) (T, error) {
	if f.wasFetched {
		return f.value, f.err
	}
	f.wasFetched = true
	if f.get == nil {
		f.err = zerrors.ThrowInternal(nil, "DOM-3gcfDV", "no getter function defined")
		return f.value, f.err
	}
	f.value, f.err = f.get(ctx, opts)
	return f.value, f.err
}

func (f *lazyGetter[T]) reload(ctx context.Context, opts *InvokeOpts) (T, error) {
	f.wasFetched = false
	f.err = nil
	var zero T
	f.value = zero
	return f.fetch(ctx, opts)
}

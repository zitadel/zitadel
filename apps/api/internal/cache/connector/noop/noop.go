package noop

import (
	"context"

	"github.com/zitadel/zitadel/internal/cache"
)

type noop[I, K comparable, V cache.Entry[I, K]] struct{}

// NewCache returns a cache that does nothing
func NewCache[I, K comparable, V cache.Entry[I, K]]() cache.Cache[I, K, V] {
	return noop[I, K, V]{}
}

func (noop[I, K, V]) Set(context.Context, V)                          {}
func (noop[I, K, V]) Get(context.Context, I, K) (value V, ok bool)    { return }
func (noop[I, K, V]) Invalidate(context.Context, I, ...K) (err error) { return }
func (noop[I, K, V]) Delete(context.Context, I, ...K) (err error)     { return }
func (noop[I, K, V]) Prune(context.Context) (err error)               { return }
func (noop[I, K, V]) Truncate(context.Context) (err error)            { return }

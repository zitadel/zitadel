package gomap

import (
	"context"
	"maps"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zitadel/zitadel/internal/cache"
)

type mapCache[I, K comparable, V any] struct {
	config         *cache.Config
	indexMap       map[I]*index[K, V]
	closeAutoPrune func()
}

// NewCache returns an in-memory Cache implementation based on the builtin go map type.
// Object values are stored as-is and there is no encoding or decoding involved.
func NewCache[I, K comparable, V any](background context.Context, indices []I, config cache.Config) cache.Cache[I, K, V] {
	m := &mapCache[I, K, V]{
		config:   &config,
		indexMap: make(map[I]*index[K, V], len(indices)),
	}
	for _, name := range indices {
		m.indexMap[name] = &index[K, V]{
			config:  m.config,
			entries: make(map[K]*entry[V]),
		}
	}
	m.closeAutoPrune = config.StartAutoPrune(background, m)
	return m
}

func (c *mapCache[I, K, V]) Get(_ context.Context, index I, key K) (value V, err error) {
	i, ok := c.indexMap[index]
	if !ok {
		return value, cache.NewIndexUnknownErr(index)
	}
	entry, err := i.Get(key)
	if err != nil {
		return value, err
	}
	return entry.value, nil
}

func (c *mapCache[I, K, V]) Set(_ context.Context, ce cache.Entry[I, K, V]) error {
	entry := &entry[V]{
		value:   ce.Value(),
		created: time.Now(),
	}
	for name, i := range c.indexMap {
		keys := ce.Keys(name)
		i.Invalidate(keys)
		i.Set(keys, entry)
	}
	return nil
}

func (c *mapCache[I, K, V]) Invalidate(_ context.Context, index I, key ...K) error {
	i, ok := c.indexMap[index]
	if !ok {
		return cache.NewIndexUnknownErr(index)
	}
	i.Invalidate(key)
	return nil
}

func (c *mapCache[I, K, V]) Delete(ctx context.Context, index I, key ...K) error {
	i, ok := c.indexMap[index]
	if !ok {
		return cache.NewIndexUnknownErr(index)
	}
	i.Delete(key)
	return nil
}

func (c *mapCache[I, K, V]) Prune(_ context.Context) error {
	for _, index := range c.indexMap {
		index.Prune()
	}
	return nil
}

func (c *mapCache[I, K, V]) Clear(_ context.Context) error {
	for _, index := range c.indexMap {
		index.Clear()
	}
	return nil
}

func (c *mapCache[I, K, V]) Close(ctx context.Context) error {
	c.closeAutoPrune()
	return ctx.Err()
}

type index[K comparable, V any] struct {
	mutex   sync.RWMutex
	config  *cache.Config
	entries map[K]*entry[V]
}

func (i *index[K, V]) Get(key K) (*entry[V], error) {
	i.mutex.RLock()
	entry, ok := i.entries[key]
	i.mutex.RUnlock()
	if ok && entry.isValid(i.config) {
		return entry, nil
	}
	return nil, cache.ErrCacheMiss
}

func (c *index[K, V]) Set(keys []K, entry *entry[V]) {
	c.mutex.Lock()
	for _, key := range keys {
		c.entries[key] = entry
	}
	c.mutex.Unlock()
}

func (i *index[K, V]) Invalidate(keys []K) {
	i.mutex.RLock()
	for _, key := range keys {
		if entry, ok := i.entries[key]; ok {
			entry.invalid.Store(true)
		}
	}
	i.mutex.RUnlock()
}

func (c *index[K, V]) Delete(keys []K) {
	c.mutex.Lock()
	for _, key := range keys {
		delete(c.entries, key)
	}
	c.mutex.Unlock()
}

func (c *index[K, V]) Prune() {
	c.mutex.Lock()
	maps.DeleteFunc(c.entries, func(_ K, entry *entry[V]) bool {
		return entry.isValid(c.config)
	})
	c.mutex.Unlock()
}

func (c *index[K, V]) Clear() {
	c.mutex.Lock()
	c.entries = make(map[K]*entry[V])
	c.mutex.Unlock()
}

type entry[V any] struct {
	value   V
	created time.Time
	invalid atomic.Bool
	lastUse atomic.Int64 // UnixMicro time
}

func (e *entry[V]) isValid(c *cache.Config) bool {
	if e.invalid.Load() {
		return false
	}
	now := time.Now()
	if c.MaxAge > 0 {
		if e.created.Add(c.MaxAge).Before(now) {
			e.invalid.Store(true)
			return false
		}
	}
	if c.LastUseAge > 0 {
		lastUse := e.lastUse.Load()
		if time.UnixMicro(lastUse).Add(c.LastUseAge).Before(now) {
			e.invalid.Store(true)
			return false
		}
		e.lastUse.CompareAndSwap(lastUse, now.UnixMicro())
	}
	return true
}

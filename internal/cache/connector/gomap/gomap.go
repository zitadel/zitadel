package gomap

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/zitadel/zitadel/internal/cache"
)

type mapCache[I, K comparable, V cache.Entry[I, K]] struct {
	config   *cache.Config
	indexMap map[I]*index[K, V]
	logger   *slog.Logger
}

// NewCache returns an in-memory Cache implementation based on the builtin go map type.
// Object values are stored as-is and there is no encoding or decoding involved.
func NewCache[I, K comparable, V cache.Entry[I, K]](background context.Context, indices []I, config cache.Config) cache.PrunerCache[I, K, V] {
	m := &mapCache[I, K, V]{
		config:   &config,
		indexMap: make(map[I]*index[K, V], len(indices)),
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelError,
		})),
	}
	if config.Log != nil {
		m.logger = config.Log.Slog()
	}
	m.logger.InfoContext(background, "map cache logging enabled")

	for _, name := range indices {
		m.indexMap[name] = &index[K, V]{
			config:  m.config,
			entries: make(map[K]*entry[V]),
		}
	}
	return m
}

func (c *mapCache[I, K, V]) Get(ctx context.Context, index I, key K) (value V, ok bool) {
	i, ok := c.indexMap[index]
	if !ok {
		c.logger.ErrorContext(ctx, "map cache get", "err", cache.NewIndexUnknownErr(index), "index", index, "key", key)
		return value, false
	}
	entry, err := i.Get(key)
	if err == nil {
		c.logger.DebugContext(ctx, "map cache get", "index", index, "key", key)
		return entry.value, true
	}
	if errors.Is(err, cache.ErrCacheMiss) {
		c.logger.InfoContext(ctx, "map cache get", "err", err, "index", index, "key", key)
		return value, false
	}
	c.logger.ErrorContext(ctx, "map cache get", "err", cache.NewIndexUnknownErr(index), "index", index, "key", key)
	return value, false
}

func (c *mapCache[I, K, V]) Set(ctx context.Context, value V) {
	now := time.Now()
	entry := &entry[V]{
		value:   value,
		created: now,
	}
	entry.lastUse.Store(now.UnixMicro())

	for name, i := range c.indexMap {
		keys := value.Keys(name)
		i.Set(keys, entry)
		c.logger.DebugContext(ctx, "map cache set", "index", name, "keys", keys)
	}
}

func (c *mapCache[I, K, V]) Invalidate(ctx context.Context, index I, keys ...K) error {
	i, ok := c.indexMap[index]
	if !ok {
		return cache.NewIndexUnknownErr(index)
	}
	i.Invalidate(keys)
	c.logger.DebugContext(ctx, "map cache invalidate", "index", index, "keys", keys)
	return nil
}

func (c *mapCache[I, K, V]) Delete(ctx context.Context, index I, keys ...K) error {
	i, ok := c.indexMap[index]
	if !ok {
		return cache.NewIndexUnknownErr(index)
	}
	i.Delete(keys)
	c.logger.DebugContext(ctx, "map cache delete", "index", index, "keys", keys)
	return nil
}

func (c *mapCache[I, K, V]) Prune(ctx context.Context) error {
	for name, index := range c.indexMap {
		index.Prune()
		c.logger.DebugContext(ctx, "map cache prune", "index", name)
	}
	return nil
}

func (c *mapCache[I, K, V]) Truncate(ctx context.Context) error {
	for name, index := range c.indexMap {
		index.Truncate()
		c.logger.DebugContext(ctx, "map cache truncate", "index", name)
	}
	return nil
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
		return !entry.isValid(c.config)
	})
	c.mutex.Unlock()
}

func (c *index[K, V]) Truncate() {
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

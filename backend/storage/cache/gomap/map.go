package gomap

import (
	"sync"

	"github.com/zitadel/zitadel/backend/storage/cache"
)

type Map[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

// Clear implements cache.Cache.
func (m *Map[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[K]V, len(m.items))
}

// Delete implements cache.Cache.
func (m *Map[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
}

// Get implements cache.Cache.
func (m *Map[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.items[key]
	return value, exists
}

// Set implements cache.Cache.
func (m *Map[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = value
}

var _ cache.Cache[string, string] = &Map[string, string]{}

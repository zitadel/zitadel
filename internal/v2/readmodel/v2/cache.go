package readmodel

import "sync"

type object interface {
}

type Cache[K comparable, V object] interface {
	Get(key K) (V, bool)
	Set(key K, value V) error
	Remove(key K) error
}

var _ Cache[string, any] = (*MapCache[string, any])(nil)

type MapCache[K comparable, V object] sync.Map

func NewMapCache[K comparable, V object]() *MapCache[K, V] {
	return new(MapCache[K, V])
}

// Get implements Cache.
func (m *MapCache[K, V]) Get(key K) (V, bool) {
	value, ok := (*sync.Map)(m).Load(key)
	if !ok {
		var v V
		return v, false
	}
	return value.(V), ok
}

// Set implements Cache.
func (m *MapCache[K, V]) Set(key K, value V) error {
	(*sync.Map)(m).Store(key, value)
	return nil
}

// Remove implements Cache.
func (m *MapCache[K, V]) Remove(key K) error {
	(*sync.Map)(m).Delete(key)
	return nil
}

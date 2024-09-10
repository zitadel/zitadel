package readmodel

type object interface {
}

type Cache[K comparable, V object] interface {
	Get(key K) (V, bool)
	Set(key K, value V) error
	Remove(key K) error
}

var _ Cache[string, any] = (*MapCache[string, any])(nil)

type MapCache[K comparable, V object] map[K]V

func NewMapCache[K comparable, V object]() *MapCache[K, V] {
	return &MapCache[K, V]{}
}

// Get implements Cache.
func (m *MapCache[K, V]) Get(key K) (V, bool) {
	value, ok := (*m)[key]
	return value, ok
}

// Set implements Cache.
func (m *MapCache[K, V]) Set(key K, value V) error {
	(*m)[key] = value
	return nil
}

// Remove implements Cache.
func (m *MapCache[K, V]) Remove(key K) error {
	delete(*m, key)
	return nil
}

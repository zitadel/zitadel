package builder

import "sync"

type builder interface {
	reset()
}

type Builder[T builder] struct {
	*sync.Pool
}

func NewBuilder[T builder](creator func() T) *Builder[T] {
	if creator == nil {
		creator = func() T {
			var x T
			return x
		}
	}
	return &Builder[T]{
		Pool: &sync.Pool{
			New: func() any {
				return creator()
			},
		},
	}
}

func (b *Builder[T]) Get() T {
	return b.Pool.Get().(T)
}

func (b *Builder[T]) Put(x T) {
	b.Pool.Put(x)
}

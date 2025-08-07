package cache

import (
	"errors"
	"fmt"
)

type IndexUnknownError[I comparable] struct {
	index I
}

func NewIndexUnknownErr[I comparable](index I) error {
	return IndexUnknownError[I]{index}
}

func (i IndexUnknownError[I]) Error() string {
	return fmt.Sprintf("index %v unknown", i.index)
}

func (a IndexUnknownError[I]) Is(err error) bool {
	if b, ok := err.(IndexUnknownError[I]); ok {
		return a.index == b.index
	}
	return false
}

var (
	ErrCacheMiss = errors.New("cache miss")
)

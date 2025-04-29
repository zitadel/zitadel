package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/receiver/cache"
)

type setCache[I, K comparable, V cache.Entry[I, K]] struct {
	cache   cache.Cache[I, K, V]
	command Command
	entry   V
}

// SetCache decorates the command, if the command is executed without error it will set the cache entry.
func SetCache[I, K comparable, V cache.Entry[I, K]](cache cache.Cache[I, K, V], command Command, entry V) Command {
	return &setCache[I, K, V]{
		cache:   cache,
		command: command,
		entry:   entry,
	}
}

var _ Command = (*setCache[any, any, cache.Entry[any, any]])(nil)

// Execute implements [Command].
func (s *setCache[I, K, V]) Execute(ctx context.Context) error {
	if err := s.command.Execute(ctx); err != nil {
		return err
	}
	s.cache.Set(ctx, s.entry)
	return nil
}

// Name implements [Command].
func (s *setCache[I, K, V]) Name() string {
	return s.command.Name()
}

type deleteCache[I, K comparable, V cache.Entry[I, K]] struct {
	cache   cache.Cache[I, K, V]
	command Command
	index   I
	keys    []K
}

// DeleteCache decorates the command, if the command is executed without error it will delete the cache entry.
func DeleteCache[I, K comparable, V cache.Entry[I, K]](cache cache.Cache[I, K, V], command Command, index I, keys ...K) Command {
	return &deleteCache[I, K, V]{
		cache:   cache,
		command: command,
		index:   index,
		keys:    keys,
	}
}

var _ Command = (*deleteCache[any, any, cache.Entry[any, any]])(nil)

// Execute implements [Command].
func (s *deleteCache[I, K, V]) Execute(ctx context.Context) error {
	if err := s.command.Execute(ctx); err != nil {
		return err
	}
	return s.cache.Delete(ctx, s.index, s.keys...)
}

// Name implements [Command].
func (s *deleteCache[I, K, V]) Name() string {
	return s.command.Name()
}

type invalidateCache[I, K comparable, V cache.Entry[I, K]] struct {
	cache   cache.Cache[I, K, V]
	command Command
	index   I
	keys    []K
}

// InvalidateCache decorates the command, if the command is executed without error it will invalidate the cache entry.
func InvalidateCache[I, K comparable, V cache.Entry[I, K]](cache cache.Cache[I, K, V], command Command, index I, keys ...K) Command {
	return &invalidateCache[I, K, V]{
		cache:   cache,
		command: command,
		index:   index,
		keys:    keys,
	}
}

var _ Command = (*invalidateCache[any, any, cache.Entry[any, any]])(nil)

// Execute implements [Command].
func (s *invalidateCache[I, K, V]) Execute(ctx context.Context) error {
	if err := s.command.Execute(ctx); err != nil {
		return err
	}
	return s.cache.Invalidate(ctx, s.index, s.keys...)
}

// Name implements [Command].
func (s *invalidateCache[I, K, V]) Name() string {
	return s.command.Name()
}

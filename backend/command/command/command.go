package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/command/receiver/cache"
)

type Command interface {
	Execute(context.Context) error
	Name() string
}

type Batch struct {
	commands []Command
}

func (b *Batch) Execute(ctx context.Context) error {
	for _, command := range b.commands {
		if err := command.Execute(ctx); err != nil {
			// TODO: undo?
			return err
		}
	}
	return nil
}

type CacheableCommand[I, K comparable, V cache.Entry[I, K]] interface {
	Command
	Entry() V
}

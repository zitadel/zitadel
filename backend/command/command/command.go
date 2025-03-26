package command

import "context"

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

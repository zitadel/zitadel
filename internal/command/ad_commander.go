package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
)

func NewCommander(agg *eventstore.Aggregate) *commander {
	return &commander{
		agg: agg,
	}
}

type commander struct {
	err             error
	agg             *eventstore.Aggregate
	shouldOverwrite bool
	previous        *commander
	command         createCommands
}

type createCommands func(context.Context, *eventstore.Aggregate) ([]eventstore.Command, error)

func (c *commander) Error() error {
	if c.err != nil || c.previous == nil {
		return c.err
	}
	return c.previous.Error()
}

func (c *commander) NextAggregate(agg *eventstore.Aggregate) *commander {
	return c.Next(nil, WithAggregate(agg), withShouldOverwrite())
}

func (c *commander) Next(cmd createCommands, opts ...commanderOption) (next *commander) {
	if c.err != nil {
		return c
	}

	if c.shouldOverwrite {
		next.shouldOverwrite = false
		next = c
	} else {
		next = &commander{
			agg:      c.agg,
			previous: c,
		}
	}
	next.command = cmd

	for _, opt := range opts {
		opt(next)
	}
	return next
}

//TODO: should be outside
// func (c *commander) exec(ctx context.Context) error {
// 	if c.err != nil {
// 		return c.err
// 	}

// 	cmds, err := c.commands(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = c.pusher.Push(ctx, cmds)

// 	return err
// }

func (c *commander) commands(ctx context.Context) ([]eventstore.Command, error) {
	if c.err != nil {
		return nil, c.err
	}
	cmds, err := c.command(ctx, c.agg)
	if err != nil || c.previous == nil {
		return cmds, err
	}

	previousCmds, err := c.previous.commands(ctx)
	if err != nil {
		return nil, err
	}

	return append(previousCmds, cmds...), nil
}

type commanderOption func(*commander)

func WithAggregate(agg *eventstore.Aggregate) commanderOption {
	return func(c *commander) {
		c.agg = agg
	}
}

func WithErr(err error) commanderOption {
	return func(c *commander) {
		c.err = err
	}
}

func withShouldOverwrite() commanderOption {
	return func(c *commander) {
		c.shouldOverwrite = true
	}
}

package command

import (
	"context"
	"errors"

	"github.com/caos/zitadel/internal/eventstore"
)

func NewCommander(agg *eventstore.Aggregate, opts ...commanderOption) *commander {
	c := new(commander)
	return c.use(append([]commanderOption{WithAggregate(agg)}, opts...))
}

type commander struct {
	err      error
	agg      *eventstore.Aggregate
	previous *commander
	command  createCommands
}

type createCommands func(context.Context, *eventstore.Aggregate) ([]eventstore.Command, error)
type commanderOption func(*commander)

var (
	ErrNotExecutable = errors.New("commander ist not executable")
	ErrNoAggregate   = errors.New("no aggregate provided")
)

func (c *commander) Next(opts ...commanderOption) (next *commander) {
	if c.err != nil {
		return c
	}

	next = &commander{
		agg:      c.agg,
		previous: c,
	}
	return next.use(opts)
}

func WithAggregate(agg *eventstore.Aggregate) commanderOption {
	return func(c *commander) {
		c.agg = agg
	}
}

func (c *commander) Commands(ctx context.Context) ([]eventstore.Command, error) {
	if c.err != nil {
		return nil, c.err
	}
	cmds, err := c.command(ctx, c.agg)
	if err != nil || c.previous == nil {
		return cmds, err
	}

	previousCmds, err := c.previous.Commands(ctx)
	if err != nil {
		return nil, err
	}

	return append(previousCmds, cmds...), nil
}

func (c *commander) Error() error {
	if c.err != nil || c.previous == nil {
		return c.err
	}
	return c.previous.Error()
}

func (c *commander) use(opts []commanderOption) *commander {
	for _, opt := range opts {
		opt(c)
	}

	if c.command == nil && c.err == nil {
		c.err = ErrNotExecutable
	}
	if c.agg == nil {
		c.err = ErrNoAggregate
	}

	return c
}

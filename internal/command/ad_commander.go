package command

import (
	"context"
	"errors"

	"github.com/caos/zitadel/internal/eventstore"
)

func NewCommander(agg *eventstore.Aggregate, filter FilterToQueryReducer, opts ...commanderOption) (*commander, error) {
	c := new(commander)
	c.defaultFilter = filter
	c.agg = agg
	return c.use(opts)
}

type commander struct {
	err           error
	agg           *eventstore.Aggregate
	previous      *commander
	commands      []createCommands
	events        []eventstore.Command
	defaultFilter FilterToQueryReducer
}

func (c *commander) filter(ctx context.Context, i *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
	events, err := c.defaultFilter(ctx, i)
	if err != nil {
		return nil, err
	}
	seq := uint64(0)
	if len(events) > 0 {
		seq = events[len(events)-1].Sequence()
	}
	for _, command := range c.events {
		seq++
		events = append(events, commandToEvent(command, seq))
	}
	return events, nil
}

func commandToEvent(command eventstore.Command, seq uint64) eventstore.Event {
	return command.(eventstore.Event)
}

type createCommands func(context.Context) ([]eventstore.Command, error)
type commanderOption func(FilterToQueryReducer) (createCommands, error)

var (
	ErrNotExecutable = errors.New("commander ist not executable")
	ErrNoAggregate   = errors.New("no aggregate provided")
)

func (c *commander) Commands(ctx context.Context) (cmds []eventstore.Command, err error) {
	if c.err != nil {
		return nil, c.err
	}
	for _, command := range c.commands {
		cmd, err := command(ctx)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd...)
		c.events = append(c.events, cmd...)
	}

	return cmds, nil
}

func (c *commander) Error() error {
	if c.err != nil || c.previous == nil {
		return c.err
	}
	return c.previous.Error()
}

func (c *commander) use(opts []commanderOption) (*commander, error) {
	for _, option := range opts {
		cmds, err := option(c.filter)
		if err != nil {
			return nil, err
		}
		c.commands = append(c.commands, cmds)
	}

	if c.commands == nil && c.err == nil {
		c.err = ErrNotExecutable
	}
	if c.agg == nil {
		c.err = ErrNoAggregate
	}

	return c, nil
}

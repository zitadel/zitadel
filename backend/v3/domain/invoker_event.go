package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
)

// eventStoreInvoker checks if the [EventProducer].Events function returns any events.
// If it does, it collects the events and publishes them to the event store.
type eventStoreInvoker struct {
	invoker
	collector eventCollector
}

func NewEventStoreInvoker(next Invoker) *eventStoreInvoker {
	return &eventStoreInvoker{
		invoker: invoker{next: next},
	}
}

type EventProducer interface {
	// Events returns the events that should be pushed to the event store after the command is executed.
	// If the command does not produce events, it should `return nil, nil`.
	Events(ctx context.Context, opts *InvokeOpts) ([]legacy_es.Command, error)
}

func (i *eventStoreInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	if _, ok := executor.(*batchExecutor); ok {
		return i.collect(ctx, executor, opts)
	}
	if _, ok := executor.(EventProducer); !ok {
		return i.execute(ctx, executor, opts)
	}
	return i.collect(ctx, executor, opts)
}

func (i *eventStoreInvoker) collect(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	if i.collector == nil {
		var close func(err error) error
		close, err = i.ensureTx(ctx, opts)
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				err = close(err)
				return
			}
			commands := i.collector.commands()
			if len(commands) == 0 {
				err = close(nil)
				return
			}
			err = eventstore.Publish(ctx, opts.LegacyEventstore(), opts.DB(), commands...)
			err = close(err)
		}()
		i.collector = initEventCollector(i.invoker, nil, executor)
	} else {
		i.collector = i.collector.initSub(executor)
	}

	defer func() {
		if i.collector.parent() == nil {
			return
		}
		i.collector = i.collector.parent()
	}()
	err = i.execute(ctx, executor, opts)
	if err != nil {
		return err
	}
	return i.collector.collect(ctx, opts)
}

type eventCollector interface {
	collect(ctx context.Context, opts *InvokeOpts) error
	commands() []legacy_es.Command
	initSub(executor Executor) eventCollector
	parent() eventCollector
}

type batchCollector struct {
	invoker
	executor   Executor
	collectors []eventCollector
	previous   eventCollector
}

// collect implements [eventCollector].
func (b batchCollector) collect(ctx context.Context, opts *InvokeOpts) error {
	return nil
}

// commands implements [eventCollector].
func (b batchCollector) commands() []legacy_es.Command {
	var commands []legacy_es.Command
	for _, collector := range b.collectors {
		commands = append(commands, collector.commands()...)
	}
	return commands
}

// initSub implements [eventCollector].
func (b *batchCollector) initSub(executor Executor) eventCollector {
	b.collectors = append(b.collectors, initEventCollector(b.invoker, b, executor))
	return b.collectors[len(b.collectors)-1]
}

// parent implements [eventCollector].
func (b batchCollector) parent() eventCollector {
	return b.previous
}

var _ eventCollector = (*batchCollector)(nil)

type commandCollector struct {
	invoker
	producer EventProducer
	sub      eventCollector
	previous eventCollector
	cmds     []legacy_es.Command
}

// collect implements [eventCollector].
func (c *commandCollector) collect(ctx context.Context, opts *InvokeOpts) (err error) {
	c.cmds, err = c.producer.Events(ctx, opts)
	return err
}

// commands implements [eventCollector].
func (c commandCollector) commands() []legacy_es.Command {
	if c.sub == nil {
		return c.cmds
	}
	return append(c.cmds, c.sub.commands()...)
}

// initSub implements [eventCollector].
func (c *commandCollector) initSub(executor Executor) eventCollector {
	c.sub = initEventCollector(c.invoker, c, executor)
	return c.sub
}

// parent implements [eventCollector].
func (c commandCollector) parent() eventCollector {
	return c.previous
}

var _ eventCollector = (*commandCollector)(nil)

func initEventCollector(i invoker, parent eventCollector, executor Executor) eventCollector {
	if _, ok := executor.(*batchExecutor); ok {
		return &batchCollector{
			invoker:  i,
			executor: executor,
			previous: parent,
		}
	}
	return &commandCollector{
		invoker:  i,
		producer: executor.(EventProducer),
		previous: parent,
	}
}

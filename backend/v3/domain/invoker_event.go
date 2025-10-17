package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
)

// eventStoreInvoker checks if the [EventProducer].Events function returns any events.
// If it does, it collects the events and publishes them to the event store.
type eventStoreInvoker struct {
	next      Invoker
	collector *eventCollector
}

func newEventStoreInvoker(next Invoker) *eventStoreInvoker {
	return &eventStoreInvoker{next: next}
}

type EventProducer interface {
	// Events returns the events that should be pushed to the event store after the command is executed.
	// If the command does not produce events, it should `return nil, nil`.
	Events(ctx context.Context, opts *InvokeOpts) ([]legacy_es.Command, error)
}

func (i *eventStoreInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	if _, ok := executor.(EventProducer); !ok {
		return i.execute(ctx, executor, opts)
	}

	if i.collector != nil {
		return i.collector.Invoke(ctx, executor, opts)
	}

	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	i.collector = &eventCollector{next: i.next}

	if err = i.collector.Invoke(ctx, executor, opts); err != nil {
		return err
	}

	if len(i.collector.events) > 0 {
		err = eventstore.Publish(ctx, legacyEventstore, opts.DB, i.collector.events...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *eventStoreInvoker) execute(ctx context.Context, executor Executor, opts *InvokeOpts) error {
	if i.next != nil {
		return i.next.Invoke(ctx, executor, opts)
	}
	return executor.Execute(ctx, opts)
}

// eventCollector collects events from all commands. The [eventStoreInvoker] pushes the collected events after all commands are executed.
// The events are collected after the command got executed, the collector ensures that the command is executed in the same transaction as writing the events.
type eventCollector struct {
	next          Invoker
	events        []legacy_es.Command
	shouldPrepend bool
}

func (i *eventCollector) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	command, ok := executor.(EventProducer)
	if !ok {
		if i.next != nil {
			return i.next.Invoke(ctx, executor, opts)
		}
		return executor.Execute(ctx, opts)
	}

	shouldPrepend := i.shouldPrepend
	i.shouldPrepend = false

	if i.next != nil {
		i.shouldPrepend = true
		err = i.next.Invoke(ctx, executor, opts)
		i.shouldPrepend = false
	} else {
		err = executor.Execute(ctx, opts)
	}
	if err != nil {
		return err
	}
	collectedEvents, err := command.Events(ctx, opts)
	if err != nil {
		return err
	}

	if shouldPrepend {
		i.events = append(collectedEvents, i.events...)
	} else {
		i.events = append(i.events, collectedEvents...)
	}

	return err
}

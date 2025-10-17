package domain

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
)

// eventStoreInvoker checks if the [Commander].Events function returns any events.
// If it does, it collects the events and publishes them to the event store.
type eventStoreInvoker struct {
	collector *eventCollector
}

func newEventStoreInvoker(next Invoker) *eventStoreInvoker {
	return &eventStoreInvoker{collector: &eventCollector{next: next}}
}

type EventProducer interface {
	// Events returns the events that should be pushed to the event store after the command is executed.
	// If the command does not produce events, it should return nil or an empty slice.
	Events(ctx context.Context, opts *InvokeOpts) ([]legacy_es.Command, error)
}

func (i *eventStoreInvoker) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	defer func() { err = close(ctx, err) }()

	err = i.collector.Invoke(ctx, executor, opts)
	if err != nil {
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

// eventCollector collects events from all commands. The [eventStoreInvoker] pushes the collected events after all commands are executed.
// The events are collected after the command got executed, the collector ensures that the command is executed in the same transaction as writing the events.
type eventCollector struct {
	next   Invoker
	events []legacy_es.Command
}

func (i *eventCollector) Invoke(ctx context.Context, executor Executor, opts *InvokeOpts) (err error) {
	command, ok := executor.(EventProducer)
	if !ok {
		if i.next != nil {
			return i.next.Invoke(ctx, executor, opts)
		}
		return executor.Execute(ctx, opts)
	}

	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	if i.next != nil {
		err = i.next.Invoke(ctx, executor, opts)
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

	i.events = append(collectedEvents, i.events...)

	return err
}

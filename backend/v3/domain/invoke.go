package domain

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	legacy_es "github.com/zitadel/zitadel/internal/eventstore"
)

type InvokeOpt func(*CommandOpts)

func WithOrganizationRepo(repo func(client database.QueryExecutor) OrganizationRepository) InvokeOpt {
	return func(opts *CommandOpts) {
		opts.organizationRepo = repo
	}
}

// Invoke provides a way to execute commands within the domain package.
// It uses a chain of responsibility pattern to handle the command execution.
// The default chain includes logging, tracing, and event publishing.
// If you want to invoke multiple commands in a single transaction, you can use the [commandBatch].
func Invoke(ctx context.Context, cmd Commander, opts ...InvokeOpt) error {
	invoker := newEventStoreInvoker(
		newLoggingInvoker(
			newTraceInvoker(
				newValidatorInvoker(nil),
			),
		),
	)
	commandOpts := &CommandOpts{
		Invoker: invoker.collector,
		DB:      pool,
	}
	for _, opt := range opts {
		opt(commandOpts)
	}
	return invoker.Invoke(ctx, cmd, commandOpts)
}

// eventStoreInvoker checks if the [Commander].Events function returns any events.
// If it does, it collects the events and publishes them to the event store.
type eventStoreInvoker struct {
	collector *eventCollector
}

func newEventStoreInvoker(next Invoker) *eventStoreInvoker {
	return &eventStoreInvoker{collector: &eventCollector{next: next}}
}

func (i *eventStoreInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	defer func() { err = close(ctx, err) }()

	err = i.collector.Invoke(ctx, command, opts)
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

func (i *eventCollector) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
	close, err := opts.EnsureTx(ctx)
	if err != nil {
		return err
	}
	defer func() { err = close(ctx, err) }()

	if i.next != nil {
		err = i.next.Invoke(ctx, command, opts)
	} else {
		err = command.Execute(ctx, opts)
	}
	if err != nil {
		return err
	}
	i.events = append(command.Events(ctx, opts), i.events...)

	return
}

// traceInvoker decorates each command with tracing.
type traceInvoker struct {
	next Invoker
}

func newTraceInvoker(next Invoker) *traceInvoker {
	return &traceInvoker{next: next}
}

func (i *traceInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%T", command))
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()

	if i.next != nil {
		return i.next.Invoke(ctx, command, opts)
	}
	return command.Execute(ctx, opts)
}

// loggingInvoker decorates each command with logging.
// It is an example implementation and logs the command name at the beginning and success or failure after the command got executed.
type loggingInvoker struct {
	next Invoker
}

func newLoggingInvoker(next Invoker) *loggingInvoker {
	return &loggingInvoker{next: next}
}

func (i *loggingInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
	logger.InfoContext(ctx, "Invoking command", "command", command.String())

	if i.next != nil {
		err = i.next.Invoke(ctx, command, opts)
	} else {
		err = command.Execute(ctx, opts)
	}

	if err != nil {
		logger.ErrorContext(ctx, "Command invocation failed", "command", command.String(), "error", err)
		return err
	}
	logger.InfoContext(ctx, "Command invocation succeeded", "command", command.String())
	return nil
}

type noopInvoker struct {
	next Invoker
}

func (i *noopInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) error {
	if i.next != nil {
		return i.next.Invoke(ctx, command, opts)
	}
	return command.Execute(ctx, opts)
}

type validatorInvoker struct {
	next Invoker
}

func newValidatorInvoker(next Invoker) *validatorInvoker {
	return &validatorInvoker{next: next}
}

func (i *validatorInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) error {
	if err := command.Validate(ctx, opts); err != nil {
		return err
	}

	if i.next != nil {
		return i.next.Invoke(ctx, command, opts)
	}

	return command.Execute(ctx, opts)
}

// cacheInvoker could be used in the future to do the caching.
// My goal would be to have two interfaces:
// - cacheSetter: which caches an object
// - cacheGetter: which gets an object from the cache, this should also skip the command execution
// type cacheInvoker struct {
// 	next Invoker
// }

// type cacher interface {
// 	Cache(opts *CommandOpts)
// }

// func (i *cacheInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
// 	if c, ok := command.(cacher); ok {
// 		c.Cache(opts)
// 	}
// 	if i.next != nil {
// 		err = i.next.Invoke(ctx, command, opts)
// 	} else {
// 		err = command.Execute(ctx, opts)
// 	}
// 	return err
// }

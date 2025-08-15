package domain

// import (
// 	"context"
// 	"fmt"

// 	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
// )

// // Invoke provides a way to execute commands within the domain package.
// // It uses a chain of responsibility pattern to handle the command execution.
// // The default chain includes logging, tracing, and event publishing.
// // If you want to invoke multiple commands in a single transaction, you can use the [commandBatch].
// func Invoke(ctx context.Context, cmd Commander) error {
// 	invoker := newEventStoreInvoker(newLoggingInvoker(newTraceInvoker(nil)))
// 	opts := &CommandOpts{
// 		Invoker: invoker.collector,
// 		DB:      pool,
// 	}
// 	return invoker.Invoke(ctx, cmd, opts)
// }

// // eventStoreInvoker checks if the command implements the [eventer] interface.
// // If it does, it collects the events and publishes them to the event store.
// type eventStoreInvoker struct {
// 	collector *eventCollector
// }

// func newEventStoreInvoker(next Invoker) *eventStoreInvoker {
// 	return &eventStoreInvoker{collector: &eventCollector{next: next}}
// }

// func (i *eventStoreInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
// 	err = i.collector.Invoke(ctx, command, opts)
// 	if err != nil {
// 		return err
// 	}
// 	if len(i.collector.events) > 0 {
// 		err = eventstore.Publish(ctx, i.collector.events, opts.DB)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// // eventCollector collects events from all commands. The [eventStoreInvoker] pushes the collected events after all commands are executed.
// type eventCollector struct {
// 	next   Invoker
// 	events []*eventstore.Event
// }

// type eventer interface {
// 	Events() []*eventstore.Event
// }

// func (i *eventCollector) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
// 	if e, ok := command.(eventer); ok && len(e.Events()) > 0 {
// 		// we need to ensure all commands are executed in the same transaction
// 		close, err := opts.EnsureTx(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		defer func() { err = close(ctx, err) }()

// 		i.events = append(i.events, e.Events()...)
// 	}
// 	if i.next != nil {
// 		return i.next.Invoke(ctx, command, opts)
// 	}
// 	return command.Execute(ctx, opts)
// }

// // traceInvoker decorates each command with tracing.
// type traceInvoker struct {
// 	next Invoker
// }

// func newTraceInvoker(next Invoker) *traceInvoker {
// 	return &traceInvoker{next: next}
// }

// func (i *traceInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
// 	ctx, span := tracer.Start(ctx, fmt.Sprintf("%T", command))
// 	defer func() {
// 		if err != nil {
// 			span.RecordError(err)
// 		}
// 		span.End()
// 	}()

// 	if i.next != nil {
// 		return i.next.Invoke(ctx, command, opts)
// 	}
// 	return command.Execute(ctx, opts)
// }

// // loggingInvoker decorates each command with logging.
// // It is an example implementation and logs the command name at the beginning and success or failure after the command got executed.
// type loggingInvoker struct {
// 	next Invoker
// }

// func newLoggingInvoker(next Invoker) *loggingInvoker {
// 	return &loggingInvoker{next: next}
// }

// func (i *loggingInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) (err error) {
// 	logger.InfoContext(ctx, "Invoking command", "command", command.String())

// 	if i.next != nil {
// 		err = i.next.Invoke(ctx, command, opts)
// 	} else {
// 		err = command.Execute(ctx, opts)
// 	}

// 	if err != nil {
// 		logger.ErrorContext(ctx, "Command invocation failed", "command", command.String(), "error", err)
// 		return err
// 	}
// 	logger.InfoContext(ctx, "Command invocation succeeded", "command", command.String())
// 	return nil
// }

// type noopInvoker struct {
// 	next Invoker
// }

// func (i *noopInvoker) Invoke(ctx context.Context, command Commander, opts *CommandOpts) error {
// 	if i.next != nil {
// 		return i.next.Invoke(ctx, command, opts)
// 	}
// 	return command.Execute(ctx, opts)
// }

// // cacheInvoker could be used in the future to do the caching.
// // My goal would be to have two interfaces:
// // - cacheSetter: which caches an object
// // - cacheGetter: which gets an object from the cache, this should also skip the command execution
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

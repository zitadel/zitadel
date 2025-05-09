package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"

	"github.com/zitadel/zitadel/backend/command/v2/pattern"
)

type command struct {
	trace.Tracer
	cmd pattern.Command
}

func Trace(tracer trace.Tracer, cmd pattern.Command) pattern.Command {
	return &command{
		Tracer: tracer,
		cmd:    cmd,
	}
}

func (cmd *command) Name() string {
	return cmd.cmd.Name()
}

func (cmd *command) Execute(ctx context.Context) error {
	ctx, span := cmd.Tracer.Start(ctx, cmd.Name())
	defer span.End()

	err := cmd.cmd.Execute(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

type query[T any] struct {
	command
	query pattern.Query[T]
}

func Query[T any](tracer trace.Tracer, q pattern.Query[T]) pattern.Query[T] {
	return &query[T]{
		command: command{
			Tracer: tracer,
			cmd:    q,
		},
		query: q,
	}
}

func (q *query[T]) Result() T {
	return q.query.Result()
}

package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type trace struct {
	command Command
	tracer  *tracing.Tracer
}

// Trace decorates the commands execute method with tracing.
// It creates a span with the command name and records any errors that occur during execution.
// The span is ended after the command is executed.
func Trace(tracer *tracing.Tracer, command Command) Command {
	return &trace{
		command: command,
		tracer:  tracer,
	}
}

// Name implements [Command].
func (l *trace) Name() string {
	return l.command.Name()
}

func (t *trace) Execute(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, t.command.Name())
	defer span.End()
	err := t.command.Execute(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

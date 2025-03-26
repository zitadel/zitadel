package command

import (
	"context"

	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type Trace struct {
	command Command
	tracer  *tracing.Tracer
}

func (t *Trace) Execute(ctx context.Context) error {
	ctx, span := t.tracer.Start(ctx, t.command.Name())
	defer span.End()
	err := t.command.Execute(ctx)
	if err != nil {
		span.RecordError(err)
	}
	return err
}

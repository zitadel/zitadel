package traced

import (
	"context"

	"github.com/zitadel/zitadel/backend/domain/factory"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type Tracer[Req, Res any] struct {
	tracing.Tracer
	next factory.Handler[Req, Res]
}

func (*Tracer[Req, Res]) Name() string {
	return "Tracer"
}

// Handle implements [factory.Handler].
func (t *Tracer[Req, Res]) Handle(ctx context.Context, request Req) (res Res, err error) {
	if t.next == nil {
		return res, nil
	}
	ctx, span := t.Tracer.Start(
		ctx,
		t.next.Name(),
	)
	defer func() {
		if err != nil {
			span.RecordError(err)
		}
		span.End()
	}()
	return t.next.Handle(ctx, request)
}

// SetNext implements [factory.Handler].
func (t *Tracer[Req, Res]) SetNext(next factory.Handler[Req, Res]) {
	t.next = next
}

// New implements [factory.Middleware].
func (t *Tracer[Req, Res]) New() factory.Handler[Req, Res] {
	return t.NewWithNext(nil)
}

// NewWithNext implements [factory.Middleware].
func (t *Tracer[Req, Res]) NewWithNext(next factory.Handler[Req, Res]) factory.Handler[Req, Res] {
	return &Tracer[Req, Res]{Tracer: t.Tracer, next: next}
}

var (
	_ factory.Middleware[any, any] = (*Tracer[any, any])(nil)
	_ factory.Handler[any, any]    = (*Tracer[any, any])(nil)
)

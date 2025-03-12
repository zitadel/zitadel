package factory

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
)

//	type Middleware[O any, H Handler[O]] interface {
//		New() H
//		NewWithNext(next Handler[O]) H
//	}
type Middleware[Req, Res any] interface {
	New() Handler[Req, Res]
	NewWithNext(next Handler[Req, Res]) Handler[Req, Res]
}

type Handler[Req, Res any] interface {
	Handle(ctx context.Context, request Req) (Res, error)
	SetNext(next Handler[Req, Res])

	Name() string
}

// type InstanceBuilder struct {
// 	tracer *traced.Instance
// 	logger *logged.Instance
// 	cache  *cache.Instance
// 	events *event.Instance
// 	db     *sql.Instance
// }

type InstanceSetUpBuilder struct {
	tracer Middleware[*repository.Instance, *repository.Instance]
	logger Middleware[*repository.Instance, *repository.Instance]
	cache  Middleware[*repository.Instance, *repository.Instance]
	events Middleware[*repository.Instance, *repository.Instance]
	db     Middleware[*repository.Instance, *repository.Instance]
}

func (i *InstanceSetUpBuilder) Build() {
	instance := i.tracer.NewWithNext(
		i.logger.NewWithNext(
			i.db.NewWithNext(
				i.events.NewWithNext(
					i.cache.New(),
				),
			),
		),
	)
	_ = instance
	// instance.
}

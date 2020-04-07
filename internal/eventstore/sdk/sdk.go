package sdk

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type filterFunc func(context.Context, *es_models.SearchQuery) ([]*es_models.Event, error)
type appendFunc func(...*es_models.Event) error
type aggregateFunc func(context.Context) (*es_models.Aggregate, error)
type pushFunc func(context.Context, ...*es_models.Aggregate) error

func Filter(ctx context.Context, filter filterFunc, appender appendFunc, query *es_models.SearchQuery) error {
	events, err := filter(ctx, query)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return errors.ThrowNotFound(nil, "EVENT-8due3", "no events found")
	}
	return appender(events...)
}

func Save(ctx context.Context, push pushFunc, aggregater aggregateFunc, appender appendFunc) error {
	aggregate, err := aggregater(ctx)
	if err != nil {
		return err
	}
	err = push(ctx, aggregate)
	if err != nil {
		return err
	}

	//error is ignored because it would be confusing if events are saved and this method would return an error
	appender(aggregate.Events...)
	return nil
}

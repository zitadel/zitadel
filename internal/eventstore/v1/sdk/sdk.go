package sdk

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type filterFunc func(context.Context, *es_models.SearchQuery) ([]*es_models.Event, error)
type appendFunc func(...*es_models.Event) error
type AggregateFunc func(context.Context) (*es_models.Aggregate, error)
type pushFunc func(context.Context, ...*es_models.Aggregate) error

func Filter(ctx context.Context, filter filterFunc, appender appendFunc, query *es_models.SearchQuery) error {
	events, err := filter(ctx, query)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return errors.ThrowNotFound(nil, "EVENT-8due3", "no events found")
	}
	err = appender(events...)
	if err != nil {
		return ThrowAppendEventError(err, "SDK-awiWK", "Errors.Internal")
	}
	return nil
}

// Push is Deprecated use PushAggregates
// Push creates the aggregates from aggregater
// and pushes the aggregates to the given pushFunc
// the given events are appended by the appender
func Push(ctx context.Context, push pushFunc, appender appendFunc, aggregaters ...AggregateFunc) (err error) {
	if len(aggregaters) < 1 {
		return errors.ThrowPreconditionFailed(nil, "SDK-q9wjp", "Errors.Internal")
	}

	aggregates, err := makeAggregates(ctx, aggregaters)
	if err != nil {
		return err
	}

	err = push(ctx, aggregates...)
	if err != nil {
		return err
	}

	return appendAggregates(appender, aggregates)
}

func PushAggregates(ctx context.Context, push pushFunc, appender appendFunc, aggregates ...*es_models.Aggregate) (err error) {
	if len(aggregates) < 1 {
		return errors.ThrowPreconditionFailed(nil, "SDK-q9wjp", "Errors.Internal")
	}

	err = push(ctx, aggregates...)
	if err != nil {
		return err
	}
	if appender == nil {
		return nil
	}

	return appendAggregates(appender, aggregates)
}

func appendAggregates(appender appendFunc, aggregates []*es_models.Aggregate) error {
	for _, aggregate := range aggregates {
		err := appender(aggregate.Events...)
		if err != nil {
			return ThrowAppendEventError(err, "SDK-o6kzK", "Errors.Internal")
		}
	}
	return nil
}

func makeAggregates(ctx context.Context, aggregaters []AggregateFunc) (aggregates []*es_models.Aggregate, err error) {
	aggregates = make([]*es_models.Aggregate, len(aggregaters))
	for i, aggregater := range aggregaters {
		aggregates[i], err = aggregater(ctx)
		if err != nil {
			return nil, err
		}
	}
	return aggregates, nil
}

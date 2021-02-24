package sdk

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
)

type filterFunc func(context.Context, *es_models.SearchQuery) ([]*es_models.Event, error)
type appendFunc func(...*es_models.Event) error
type AggregateFunc func(context.Context) (*es_models.Aggregate, error)

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

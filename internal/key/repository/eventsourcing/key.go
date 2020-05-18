package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/key/repository/eventsourcing/model"
)

func KeyPairAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, pair *model.KeyPair) (*es_models.Aggregate, error) {
	if pair == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d5HNJA", "existing key pair must not be nil")
	}
	return aggCreator.NewAggregate(ctx, pair.AggregateID, model.KeyPairAggregate, model.KeyPairVersion, pair.Sequence)
}

func KeyPairCreateAggregate(aggCreator *es_models.AggregateCreator, pair *model.KeyPair) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		agg, err := KeyPairAggregate(ctx, aggCreator, pair)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.KeyPairAdded, pair)
	}
}

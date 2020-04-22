package eventsourcing

import (
	"context"
	"strconv"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
	"github.com/sony/sonyflake"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

func PolicyQuery(latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.PolicyAggregate).
		LatestSequenceFilter(latestSequence)
}

func PolicyAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, id string, sequence uint64) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, id, model.PolicyAggregate, projectVersion, sequence)
}

func PolicyCreateAggregate(aggCreator *es_models.AggregateCreator, project *PasswordComplexityPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if project == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "project should not be nil")
		}
		var err error
		id, err := idGenerator.NextID()
		if err != nil {
			return nil, err
		}
		project.AggregateID = strconv.FormatUint(id, 10)

		agg, err := PolicyAggregate(ctx, aggCreator, project.AggregateID, project.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.PolicyAdded, project)
	}
}

func PolicyUpdateAggregate(aggCreator *es_models.AggregateCreator, existing *PasswordComplexityPolicy, new *PasswordComplexityPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
		}
		if new == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
		}
		agg, err := PolicyAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.Changes(new)
		return agg.AppendEvent(model.PolicyChanged, changes)
	}
}

package eventsourcing

import (
	"context"
	"strconv"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
)

func PasswordAgePolicyQuery(recourceOwner string, latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.PasswordAgePolicyAggregate).
		LatestSequenceFilter(latestSequence).
		ResourceOwnerFilter(recourceOwner)

}
func PasswordAgePolicyAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, id string, sequence uint64) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, id, model.PasswordAgePolicyAggregate, projectVersion, sequence)
}

func PasswordAgePolicyCreateAggregate(aggCreator *es_models.AggregateCreator, project *PasswordAgePolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
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

		agg, err := PasswordAgePolicyAggregate(ctx, aggCreator, project.AggregateID, project.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.PasswordAgePolicyAdded, project)
	}
}

func PasswordAgePolicyUpdateAggregate(aggCreator *es_models.AggregateCreator, existing *PasswordAgePolicy, new *PasswordAgePolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
		}
		if new == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
		}
		agg, err := PasswordAgePolicyAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.AgeChanges(new)
		return agg.AppendEvent(model.PasswordAgePolicyChanged, changes)
	}
}

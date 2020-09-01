package eventsourcing

import (
	"context"

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

func PasswordAgePolicyAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, policy *PasswordAgePolicy) (*es_models.Aggregate, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-1T05i", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, policy.AggregateID, model.PasswordAgePolicyAggregate, policyAgeVersion, policy.Sequence)
}

func PasswordAgePolicyCreateAggregate(aggCreator *es_models.AggregateCreator, policy *PasswordAgePolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "Errors.Internal")
		}
		agg, err := PasswordAgePolicyAggregate(ctx, aggCreator, policy)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.PasswordAgePolicyAdded, policy)
	}
}

func PasswordAgePolicyUpdateAggregate(aggCreator *es_models.AggregateCreator, existingPolicy *PasswordAgePolicy, newPolicy *PasswordAgePolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if newPolicy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
		}
		agg, err := PasswordAgePolicyAggregate(ctx, aggCreator, existingPolicy)
		if err != nil {
			return nil, err
		}
		changes := existingPolicy.AgeChanges(newPolicy)
		return agg.AppendEvent(model.PasswordAgePolicyChanged, changes)
	}
}

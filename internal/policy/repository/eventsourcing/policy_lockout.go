package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/policy/model"
)

func PasswordLockoutPolicyQuery(recourceOwner string, latestSequence uint64) *es_models.SearchQuery {
	return es_models.NewSearchQuery().
		AggregateTypeFilter(model.PasswordLockoutPolicyAggregate).
		LatestSequenceFilter(latestSequence).
		ResourceOwnerFilter(recourceOwner)

}

func PasswordLockoutPolicyAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, policy *PasswordLockoutPolicy) (*es_models.Aggregate, error) {
	if policy == nil {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-aTRlj", "Errors.Internal")
	}
	return aggCreator.NewAggregate(ctx, policy.AggregateID, model.PasswordLockoutPolicyAggregate, policyLockoutVersion, policy.Sequence)
}

func PasswordLockoutPolicyCreateAggregate(aggCreator *es_models.AggregateCreator, policy *PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "Errors.Internal")
		}

		agg, err := PasswordLockoutPolicyAggregate(ctx, aggCreator, policy)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.PasswordLockoutPolicyAdded, policy)
	}
}

func PasswordLockoutPolicyUpdateAggregate(aggCreator *es_models.AggregateCreator, existingPolicy *PasswordLockoutPolicy, newPolicy *PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if newPolicy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "Errors.Internal")
		}
		agg, err := PasswordLockoutPolicyAggregate(ctx, aggCreator, existingPolicy)
		if err != nil {
			return nil, err
		}
		changes := existingPolicy.LockoutChanges(newPolicy)
		return agg.AppendEvent(model.PasswordLockoutPolicyChanged, changes)
	}
}

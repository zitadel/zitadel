package eventsourcing

import (
	"context"
	"strconv"

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
	return aggCreator.NewAggregate(ctx, policy.AggregateID, model.PasswordLockoutPolicyAggregate, policyLockoutVersion, policy.Sequence)
}

func PasswordLockoutPolicyCreateAggregate(aggCreator *es_models.AggregateCreator, policy *PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-kdie6", "policy should not be nil")
		}
		var err error
		id, err := idGenerator.NextID()
		if err != nil {
			return nil, err
		}
		policy.AggregateID = strconv.FormatUint(id, 10)

		agg, err := PasswordLockoutPolicyAggregate(ctx, aggCreator, policy)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.PasswordLockoutPolicyAdded, policy)
	}
}

func PasswordLockoutPolicyUpdateAggregate(aggCreator *es_models.AggregateCreator, existing *PasswordLockoutPolicy, new *PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing policy should not be nil")
		}
		if new == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new policy should not be nil")
		}
		agg, err := PasswordLockoutPolicyAggregate(ctx, aggCreator, existing)
		if err != nil {
			return nil, err
		}
		changes := existing.LockoutChanges(new)
		return agg.AppendEvent(model.PasswordLockoutPolicyChanged, changes)
	}
}

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
func PasswordLockoutPolicyAggregate(ctx context.Context, aggCreator *es_models.AggregateCreator, id string, sequence uint64) (*es_models.Aggregate, error) {
	return aggCreator.NewAggregate(ctx, id, model.PasswordLockoutPolicyAggregate, projectVersion, sequence)
}

func PasswordLockoutPolicyCreateAggregate(aggCreator *es_models.AggregateCreator, project *PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
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

		agg, err := PasswordLockoutPolicyAggregate(ctx, aggCreator, project.AggregateID, project.Sequence)
		if err != nil {
			return nil, err
		}

		return agg.AppendEvent(model.PasswordLockoutPolicyAdded, project)
	}
}

func PasswordLockoutPolicyUpdateAggregate(aggCreator *es_models.AggregateCreator, existing *PasswordLockoutPolicy, new *PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dk93d", "existing project should not be nil")
		}
		if new == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-dhr74", "new project should not be nil")
		}
		agg, err := PasswordLockoutPolicyAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.LockoutChanges(new)
		return agg.AppendEvent(model.PasswordLockoutPolicyChanged, changes)
	}
}

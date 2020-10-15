package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func PasswordLockoutPolicyAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4Tfi9", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingPasswordLockoutPolicyValidation()
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.PasswordLockoutPolicyAdded, policy)
	}
}

func PasswordLockoutPolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.PasswordLockoutPolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-5lf9s", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.PasswordLockoutPolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-1dj9G", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.PasswordLockoutPolicyChanged, changes)
	}
}

func PasswordLockoutPolicyRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-6Hsl0", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.PasswordLockoutPolicyRemoved, nil)
	}
}

func checkExistingPasswordLockoutPolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existing := false
		for _, event := range events {
			switch event.Type {
			case model.PasswordLockoutPolicyAdded:
				existing = true
			case model.PasswordLockoutPolicyRemoved:
				existing = false
			}
		}
		if existing {
			return errors.ThrowPreconditionFailed(nil, "EVENT-2Dfl0", "Errors.Org.PasswordLockoutPolicy.AlreadyExists")
		}
		return nil
	}
}

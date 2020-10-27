package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func PasswordAgePolicyAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.PasswordAgePolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Hms9d", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingPasswordAgePolicyValidation()
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.PasswordAgePolicyAdded, policy)
	}
}

func PasswordAgePolicyChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.PasswordAgePolicy) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-3Gms0f", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.PasswordAgePolicy.Changes(policy)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Ls0rs", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.PasswordAgePolicyChanged, changes)
	}
}

func PasswordAgePolicyRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4Gsk9", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.PasswordAgePolicyRemoved, nil)
	}
}

func checkExistingPasswordAgePolicyValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existing := false
		for _, event := range events {
			switch event.Type {
			case model.PasswordAgePolicyAdded:
				existing = true
			case model.PasswordAgePolicyRemoved:
				existing = false
			}
		}
		if existing {
			return errors.ThrowPreconditionFailed(nil, "EVENT-2Fjs0", "Errors.Org.PasswordAgePolicy.AlreadyExists")
		}
		return nil
	}
}

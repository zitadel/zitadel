package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func MailTemplateAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, policy *iam_es_model.MailTemplate) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if policy == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-4BeRi", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingMailTemplateValidation()
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.MailTemplateAdded, policy)
	}
}

func MailTemplateChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, template *iam_es_model.MailTemplate) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if template == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-yzXO0", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := existing.MailTemplate.Changes(template)
		if len(changes) == 0 {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-erTCI", "Errors.NoChangesFound")
		}
		return agg.AppendEvent(model.MailTemplateChanged, changes)
	}
}

func MailTemplateRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-2jVit", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		return agg.AppendEvent(model.MailTemplateRemoved, nil)
	}
}

func checkExistingMailTemplateValidation() func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existing := false
		for _, event := range events {
			switch event.Type {
			case model.MailTemplateAdded:
				existing = true
			case model.MailTemplateRemoved:
				existing = false
			}
		}
		if existing {
			return errors.ThrowPreconditionFailed(nil, "EVENT-aUH4D", "Errors.Org.MailTemplate.AlreadyExists")
		}
		return nil
	}
}

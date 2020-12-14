package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

func MailTextAddedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, mailText *iam_es_model.MailText) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mailText == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Gk3Cn", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		validationQuery := es_models.NewSearchQuery().
			AggregateTypeFilter(model.OrgAggregate).
			AggregateIDFilter(existing.AggregateID)

		validation := checkExistingMailTextValidation(mailText, existing.MailTexts)
		agg.SetPrecondition(validationQuery, validation)
		return agg.AppendEvent(model.MailTextAdded, mailText)
	}
}

func MailTextChangedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, mailText *iam_es_model.MailText) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if mailText == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-Hog8a", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 2)
		for _, exMailText := range existing.MailTexts {
			if exMailText.MailTextType == mailText.MailTextType && exMailText.Language == mailText.Language {
				changes = exMailText.Changes(mailText)
				if len(changes) == 0 {
					return nil, errors.ThrowPreconditionFailed(nil, "EVENT-DuRxA", "Errors.NoChangesFound")
				}
			}
		}
		return agg.AppendEvent(model.MailTextChanged, changes)
	}
}

func MailTextRemovedAggregate(aggCreator *es_models.AggregateCreator, existing *model.Org, mailText *iam_es_model.MailText) func(ctx context.Context) (*es_models.Aggregate, error) {
	return func(ctx context.Context) (*es_models.Aggregate, error) {
		if existing == nil {
			return nil, errors.ThrowPreconditionFailed(nil, "EVENT-cJ5Wp", "Errors.Internal")
		}
		agg, err := OrgAggregate(ctx, aggCreator, existing.AggregateID, existing.Sequence)
		if err != nil {
			return nil, err
		}
		changes := make(map[string]interface{}, 2)
		for _, exMailText := range existing.MailTexts {
			if exMailText.MailTextType == mailText.MailTextType && exMailText.Language == mailText.Language {
				mailText.ButtonText = exMailText.ButtonText
				mailText.Greeting = exMailText.Greeting
				mailText.Text = exMailText.Text
				mailText.Title = exMailText.Title
				mailText.Subject = exMailText.Subject
				mailText.PreHeader = exMailText.PreHeader
				changes = exMailText.Changes(mailText)
				if len(changes) == 0 {
					return nil, errors.ThrowPreconditionFailed(nil, "EVENT-DuRxA", "Errors.NoChangesFound")
				}
			}
		}
		return agg.AppendEvent(model.MailTextRemoved, changes)
	}
}

func checkExistingMailTextValidation(mailText *iam_es_model.MailText, existingMailTexts []*iam_es_model.MailText) func(...*es_models.Event) error {
	return func(events ...*es_models.Event) error {
		existing := false
		for _, text := range existingMailTexts {
			if text.MailTextType == mailText.MailTextType && text.Language == mailText.Language {
				existing = true
			}
		}
		if existing {
			return errors.ThrowPreconditionFailed(nil, "EVENT-zEZh7", "Errors.Org.MailText.AlreadyExists")
		}
		return nil
	}
}

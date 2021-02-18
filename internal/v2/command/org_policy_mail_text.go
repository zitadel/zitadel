package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddMailText(ctx context.Context, resourceOwner string, mailText *domain.MailText) (*domain.MailText, error) {
	if !mailText.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4778u", "Errors.Org.MailText.Invalid")
	}
	addedPolicy := NewOrgMailTextWriteModel(resourceOwner, mailText.MailTextType, mailText.Language)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-9kufs", "Errors.Org.MailText.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.MailTextWriteModel.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(
		ctx,
		org.NewMailTextAddedEvent(
			ctx,
			orgAgg,
			mailText.MailTextType,
			mailText.Language,
			mailText.Title,
			mailText.PreHeader,
			mailText.Subject,
			mailText.Greeting,
			mailText.Text,
			mailText.ButtonText))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToMailText(&addedPolicy.MailTextWriteModel), nil
}

func (r *CommandSide) ChangeMailText(ctx context.Context, resourceOwner string, mailText *domain.MailText) (*domain.MailText, error) {
	if !mailText.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-3m9fs", "Errors.Org.MailText.Invalid")
	}
	existingPolicy := NewOrgMailTextWriteModel(resourceOwner, mailText.MailTextType, mailText.Language)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-3n8fM", "Errors.Org.MailText.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.MailTextWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		orgAgg,
		mailText.MailTextType,
		mailText.Language,
		mailText.Title,
		mailText.PreHeader,
		mailText.Subject,
		mailText.Greeting,
		mailText.Text,
		mailText.ButtonText)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-2n9fs", "Errors.Org.MailText.NotChanged")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToMailText(&existingPolicy.MailTextWriteModel), nil
}

func (r *CommandSide) RemoveMailText(ctx context.Context, resourceOwner, mailTextType, language string) error {
	existingPolicy := NewOrgMailTextWriteModel(resourceOwner, mailTextType, language)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3b8Jf", "Errors.Org.MailText.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, org.NewMailTextRemovedEvent(ctx, orgAgg, mailTextType, language))
	return err
}

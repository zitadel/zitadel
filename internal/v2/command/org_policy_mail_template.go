package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddMailTemplate(ctx context.Context, resourceOwner string, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-3m9fs", "Errors.Org.MailTemplate.Invalid")
	}
	addedPolicy := NewOrgMailTemplateWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-9kufs", "Errors.Org.MailTemplate.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.MailTemplateWriteModel.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(ctx, org.NewMailTemplateAddedEvent(ctx, orgAgg, policy.Template))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToMailTemplate(&addedPolicy.MailTemplateWriteModel), nil
}

func (r *CommandSide) ChangeMailTemplate(ctx context.Context, resourceOwner string, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-9f9ds", "Errors.Org.MailTemplate.Invalid")
	}
	existingPolicy := NewOrgMailTemplateWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-5m9ie", "Errors.Org.MailTemplate.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.MailTemplateWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.Template)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9vs", "Errors.Org.MailTemplate.NotChanged")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToMailTemplate(&existingPolicy.MailTemplateWriteModel), nil
}

func (r *CommandSide) RemoveMailTemplate(ctx context.Context, orgID string) error {
	existingPolicy := NewOrgMailTemplateWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3b8Jf", "Errors.Org.MailTemplate.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)

	_, err = r.eventstore.PushEvents(ctx, org.NewMailTemplateRemovedEvent(ctx, orgAgg))
	return err
}

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
	orgAgg.PushEvents(org.NewMailTemplateAddedEvent(ctx, policy.Template))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToMailTemplate(&addedPolicy.MailTemplateWriteModel), nil
}

func (r *CommandSide) ChangeMailTemplate(ctx context.Context, resourceOwner string, policy *domain.MailTemplate) (*domain.MailTemplate, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-9f9ds", "Errors.Org.MailTemplate.Invalid")
	}
	existingPolicy := NewOrgMailTemplateWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-5m9ie", "Errors.Org.MailTemplate.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.Template)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9vs", "Errors.Org.MailTemplate.NotChanged")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.MailTemplateWriteModel.WriteModel)
	orgAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
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
	orgAgg.PushEvents(org.NewMailTemplateRemovedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
}

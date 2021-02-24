package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddLabelPolicy(ctx context.Context, resourceOwner string, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	addedPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LabelPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyAddedEvent(ctx, orgAgg, policy.PrimaryColor, policy.SecondaryColor))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLabelPolicy(&addedPolicy.LabelPolicyWriteModel), nil
}

func (c *Commands) ChangeLabelPolicy(ctx context.Context, resourceOwner string, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	existingPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-0K9dq", "Errors.Org.LabelPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.PrimaryColor, policy.SecondaryColor)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9vs", "Errors.Org.LabelPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLabelPolicy(&existingPolicy.LabelPolicyWriteModel), nil
}

func (c *Commands) RemoveLabelPolicy(ctx context.Context, orgID string) error {
	existingPolicy := NewOrgLabelPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3M9df", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	_, err = c.eventstore.PushEvents(ctx, org.NewLabelPolicyRemovedEvent(ctx, orgAgg))
	return err
}

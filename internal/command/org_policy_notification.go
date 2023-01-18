package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func (c *Commands) AddNotificationPolicy(ctx context.Context, resourceOwner string, policy *domain.NotificationPolicy) (*domain.NotificationPolicy, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-x801sk2i", "Errors.ResourceOwnerMissing")
	}
	addedPolicy := NewOrgNotificationPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-xa08n2", "Errors.Org.NotificationPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	pushedEvents, err := c.eventstore.Push(
		ctx,
		org.NewNotificationPolicyAddedEvent(
			ctx,
			orgAgg,
			policy.PasswordChange,
		))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToNotificationPolicy(&addedPolicy.NotificationPolicyWriteModel), nil
}

func (c *Commands) ChangeNotificationPolicy(ctx context.Context, resourceOwner string, policy *domain.NotificationPolicy) (*domain.NotificationPolicy, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-x091n1g", "Errors.ResourceOwnerMissing")
	}

	existingPolicy := NewOrgNotificationPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-x029n3", "Errors.Org.NotificationPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.NotificationPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.PasswordChange)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-ioqnxz", "Errors.Org.NotificationPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToNotificationPolicy(&existingPolicy.NotificationPolicyWriteModel), nil
}

func (c *Commands) RemoveNotificationPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-x89ns2", "Errors.ResourceOwnerMissing")
	}
	existingPolicy := NewOrgNotificationPolicyWriteModel(orgID)
	event, err := c.removeNotificationPolicy(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.NotificationPolicyWriteModel.WriteModel), nil
}

func (c *Commands) removeNotificationPolicy(ctx context.Context, existingPolicy *OrgNotificationPolicyWriteModel) (*org.NotificationPolicyRemovedEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-x029n1s", "Errors.Org.NotificationPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewNotificationPolicyRemovedEvent(ctx, orgAgg), nil
}

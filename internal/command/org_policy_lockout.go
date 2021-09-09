package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddLockoutPolicy(ctx context.Context, resourceOwner string, policy *domain.LockoutPolicy) (*domain.LockoutPolicy, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-8fJif", "Errors.ResourceOwnerMissing")
	}
	addedPolicy, err := c.orgLockoutPolicyWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-0olDf", "Errors.ORG.LockoutPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLockoutPolicyAddedEvent(ctx, orgAgg, policy.MaxPasswordAttempts, policy.ShowLockOutFailures))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLockoutPolicy(&addedPolicy.LockoutPolicyWriteModel), nil
}

func (c *Commands) ChangeLockoutPolicy(ctx context.Context, resourceOwner string, policy *domain.LockoutPolicy) (*domain.LockoutPolicy, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-3J9fs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLockoutPolicyWriteModelByID(ctx, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-ADfs1", "Errors.Org.LockoutPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LockoutPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.MaxPasswordAttempts, policy.ShowLockOutFailures)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-4M9vs", "Errors.Org.LockoutPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLockoutPolicy(&existingPolicy.LockoutPolicyWriteModel), nil
}

func (c *Commands) RemoveLockoutPolicy(ctx context.Context, orgID string) error {
	if orgID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "Org-4J9fs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLockoutPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-D4zuz", "Errors.Org.LockoutPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)

	_, err = c.eventstore.PushEvents(ctx, org.NewLockoutPolicyRemovedEvent(ctx, orgAgg))
	return err
}

func (c *Commands) removeLockoutPolicyIfExists(ctx context.Context, orgID string) (*org.LockoutPolicyRemovedEvent, error) {
	existingPolicy, err := c.orgLockoutPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State != domain.PolicyStateActive {
		return nil, nil
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewLockoutPolicyRemovedEvent(ctx, orgAgg), nil
}

func (c *Commands) orgLockoutPolicyWriteModelByID(ctx context.Context, orgID string) (*OrgLockoutPolicyWriteModel, error) {
	policy := NewOrgLockoutPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

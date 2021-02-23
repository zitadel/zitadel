package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (r *CommandSide) AddPasswordLockoutPolicy(ctx context.Context, resourceOwner string, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	addedPolicy := NewOrgPasswordLockoutPolicyWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-0olDf", "Errors.ORG.PasswordLockoutPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(ctx, org.NewPasswordLockoutPolicyAddedEvent(ctx, orgAgg, policy.MaxAttempts, policy.ShowLockOutFailures))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordLockoutPolicy(&addedPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (r *CommandSide) ChangePasswordLockoutPolicy(ctx context.Context, resourceOwner string, policy *domain.PasswordLockoutPolicy) (*domain.PasswordLockoutPolicy, error) {
	existingPolicy := NewOrgPasswordLockoutPolicyWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-ADfs1", "Errors.Org.PasswordLockoutPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PasswordLockoutPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.MaxAttempts, policy.ShowLockOutFailures)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "ORG-4M9vs", "Errors.Org.PasswordLockoutPolicy.NotChanged")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordLockoutPolicy(&existingPolicy.PasswordLockoutPolicyWriteModel), nil
}

func (r *CommandSide) RemovePasswordLockoutPolicy(ctx context.Context, orgID string) error {
	existingPolicy := NewOrgPasswordLockoutPolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-D4zuz", "Errors.Org.PasswordLockoutPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)

	_, err = r.eventstore.PushEvents(ctx, org.NewPasswordLockoutPolicyRemovedEvent(ctx, orgAgg))
	return err
}

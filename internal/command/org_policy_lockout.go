package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddLockoutPolicy(ctx context.Context, resourceOwner string, policy *domain.LockoutPolicy) (_ *domain.LockoutPolicy, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-8fJif", "Errors.ResourceOwnerMissing")
	}
	addedPolicy, err := orgLockoutPolicyWriteModelByID(ctx, resourceOwner, c.eventstore.FilterToQueryReducer)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "ORG-0olDf", "Errors.ORG.LockoutPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLockoutPolicyAddedEvent(
		ctx,
		orgAgg,
		policy.MaxPasswordAttempts,
		policy.MaxOTPAttempts,
		policy.ShowLockOutFailures,
	))
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
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-3J9fs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := orgLockoutPolicyWriteModelByID(ctx, resourceOwner, c.eventstore.FilterToQueryReducer)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-ADfs1", "Errors.Org.LockoutPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LockoutPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.MaxPasswordAttempts, policy.MaxOTPAttempts, policy.ShowLockOutFailures)
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "ORG-0JFSr", "Errors.Org.LockoutPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLockoutPolicy(&existingPolicy.LockoutPolicyWriteModel), nil
}

func (c *Commands) RemoveLockoutPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-4J9fs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := orgLockoutPolicyWriteModelByID(ctx, orgID, c.eventstore.FilterToQueryReducer)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-D4zuz", "Errors.Org.LockoutPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)

	pushedEvents, err := c.eventstore.Push(ctx, org.NewLockoutPolicyRemovedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LockoutPolicyWriteModel.WriteModel), nil

}

func (c *Commands) removeLockoutPolicyIfExists(ctx context.Context, orgID string) (*org.LockoutPolicyRemovedEvent, error) {
	existingPolicy, err := orgLockoutPolicyWriteModelByID(ctx, orgID, c.eventstore.FilterToQueryReducer)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State != domain.PolicyStateActive {
		return nil, nil
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewLockoutPolicyRemovedEvent(ctx, orgAgg), nil
}

func orgLockoutPolicyWriteModelByID(ctx context.Context, orgID string, queryReducer func(ctx context.Context, r eventstore.QueryReducer) error) (_ *OrgLockoutPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	policy := NewOrgLockoutPolicyWriteModel(orgID)
	err = queryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func getLockoutPolicy(ctx context.Context, orgID string, queryReducer func(ctx context.Context, r eventstore.QueryReducer) error) (*domain.LockoutPolicy, error) {
	orgWm, err := orgLockoutPolicyWriteModelByID(ctx, orgID, queryReducer)
	if err != nil {
		return nil, err
	}
	if orgWm.State == domain.PolicyStateActive {
		return writeModelToLockoutPolicy(&orgWm.LockoutPolicyWriteModel), nil
	}
	instanceWm, err := defaultLockoutPolicyWriteModelByID(ctx, queryReducer)
	if err != nil {
		return nil, err
	}
	policy := writeModelToLockoutPolicy(&instanceWm.LockoutPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (r *CommandSide) AddPasswordAgePolicy(ctx context.Context, resourceOwner string, policy *domain.PasswordAgePolicy) (*domain.PasswordAgePolicy, error) {
	addedPolicy := NewOrgPasswordAgePolicyWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "ORG-Lk0dS", "Errors.Org.PasswordAgePolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.WriteModel)
	pushedEvents, err := r.eventstore.PushEvents(ctx, org.NewPasswordAgePolicyAddedEvent(ctx, orgAgg, policy.ExpireWarnDays, policy.MaxAgeDays))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordAgePolicy(&addedPolicy.PasswordAgePolicyWriteModel), nil
}

func (r *CommandSide) ChangePasswordAgePolicy(ctx context.Context, resourceOwner string, policy *domain.PasswordAgePolicy) (*domain.PasswordAgePolicy, error) {
	existingPolicy := NewOrgPasswordAgePolicyWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-0oPew", "Errors.Org.PasswordAgePolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PasswordAgePolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, orgAgg, policy.ExpireWarnDays, policy.MaxAgeDays)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-dsgjR", "Errors.ORg.LabelPolicy.NotChanged")
	}

	pushedEvents, err := r.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordAgePolicy(&existingPolicy.PasswordAgePolicyWriteModel), nil
}

func (r *CommandSide) RemovePasswordAgePolicy(ctx context.Context, orgID string) error {
	existingPolicy := NewOrgPasswordAgePolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "ORG-Dgs1g", "Errors.Org.PasswordAgePolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	_, err = r.eventstore.PushEvents(ctx, org.NewPasswordAgePolicyRemovedEvent(ctx, orgAgg))
	return err
}

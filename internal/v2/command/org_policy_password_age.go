package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
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
	orgAgg.PushEvents(org.NewPasswordAgePolicyAddedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, orgAgg)
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

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.ExpireWarnDays, policy.MaxAgeDays)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-dsgjR", "Errors.ORg.LabelPolicy.NotChanged")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.PasswordAgePolicyWriteModel.WriteModel)
	orgAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
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
	orgAgg.PushEvents(org.NewPasswordAgePolicyRemovedEvent(ctx))
	return r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
}

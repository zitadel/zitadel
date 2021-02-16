package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/org"
)

func (r *CommandSide) AddLabelPolicy(ctx context.Context, resourceOwner string, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	addedPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LabelPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	orgAgg.PushEvents(org.NewLabelPolicyAddedEvent(ctx, policy.PrimaryColor, policy.SecondaryColor))

	err = r.eventstore.PushAggregate(ctx, addedPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLabelPolicy(&addedPolicy.LabelPolicyWriteModel), nil
}

func (r *CommandSide) ChangeLabelPolicy(ctx context.Context, resourceOwner string, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	existingPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-0K9dq", "Errors.Org.LabelPolicy.NotFound")
	}

	changedEvent, hasChanged := existingPolicy.NewChangedEvent(ctx, policy.PrimaryColor, policy.SecondaryColor)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9vs", "Errors.Org.LabelPolicy.NotChanged")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	orgAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToLabelPolicy(&existingPolicy.LabelPolicyWriteModel), nil
}

func (r *CommandSide) RemoveLabelPolicy(ctx context.Context, orgID string) error {
	existingPolicy := NewOrgLabelPolicyWriteModel(orgID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return caos_errs.ThrowNotFound(nil, "Org-3M9df", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	orgAgg.PushEvents(org.NewLabelPolicyRemovedEvent(ctx))

	return r.eventstore.PushAggregate(ctx, existingPolicy, orgAgg)
}

package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) AddLabelPolicy(ctx context.Context, resourceOwner string, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Fn8ds", "Errors.ResourceOwnerMissing")
	}
	if !policy.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Md9sf", "Errors.Org.LabelPolicy.Invalid")
	}
	addedPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LabelPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyAddedEvent(
		ctx,
		orgAgg,
		policy.PrimaryColor,
		policy.SecondaryColor,
		policy.WarnColor,
		policy.PrimaryColorDark,
		policy.SecondaryColorDark,
		policy.WarnColorDark,
		policy.HideLoginNameSuffix,
		policy.ErrorMsgPopup,
		policy.DisableWatermark))
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
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-3N9fs", "Errors.ResourceOwnerMissing")
	}
	if !policy.IsValid() {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-dM9fs", "Errors.Org.LabelPolicy.Invalid")
	}
	existingPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-0K9dq", "Errors.Org.LabelPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		orgAgg,
		policy.PrimaryColor,
		policy.SecondaryColor,
		policy.WarnColor,
		policy.PrimaryColorDark,
		policy.SecondaryColorDark,
		policy.WarnColorDark,
		policy.HideLoginNameSuffix,
		policy.ErrorMsgPopup,
		policy.DisableWatermark)
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

func (c *Commands) ActivateLabelPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-KKd4X", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-34mSE", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyActivatedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddLogoLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-KKd4X", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-23BMs", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyLogoAddedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveLogoLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-2FN8s", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-4MVsf", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyLogoRemovedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddIconLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-hMDs3", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-4nq2f", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyIconAddedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIconLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-1nd0d", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-1nd9f", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyIconRemovedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddLogoDarkLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-67Ms2", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-QSqcd", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyLogoDarkAddedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveLogoDarkLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-4NF0d", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-0peQw", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyLogoDarkRemovedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddIconDarkLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-tzBfs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-4Nf8s", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyIconDarkAddedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIconDarkLabelPolicy(ctx context.Context, orgID, storageKey string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Mv9ds", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "ORG-3NFos", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewLabelPolicyIconDarkRemovedEvent(ctx, orgAgg, storageKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveLabelPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-Mf9sf", "Errors.ResourceOwnerMissing")
	}
	existingPolicy := NewOrgLabelPolicyWriteModel(orgID)
	removeEvent, err := c.removeLabelPolicy(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, removeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) removeLabelPolicy(ctx context.Context, existingPolicy *OrgLabelPolicyWriteModel) (*org.LabelPolicyRemovedEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Org-3M9df", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewLabelPolicyRemovedEvent(ctx, orgAgg), nil
}

func (c *Commands) removeLabelPolicyIfExists(ctx context.Context, orgID string) (*org.LabelPolicyRemovedEvent, error) {
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State != domain.PolicyStateActive {
		return nil, nil
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewLabelPolicyRemovedEvent(ctx, orgAgg), nil
}

func (c *Commands) orgLabelPolicyWriteModelByID(ctx context.Context, orgID string) (*OrgLabelPolicyWriteModel, error) {
	policy := NewOrgLabelPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

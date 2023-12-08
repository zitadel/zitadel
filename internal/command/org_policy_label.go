package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddLabelPolicy(ctx context.Context, resourceOwner string, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-Fn8ds", "Errors.ResourceOwnerMissing")
	}
	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	addedPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, zerrors.ThrowAlreadyExists(nil, "Org-2B0ps", "Errors.Org.LabelPolicy.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyAddedEvent(
		ctx,
		orgAgg,
		policy.PrimaryColor,
		policy.BackgroundColor,
		policy.WarnColor,
		policy.FontColor,
		policy.PrimaryColorDark,
		policy.BackgroundColorDark,
		policy.WarnColorDark,
		policy.FontColorDark,
		policy.HideLoginNameSuffix,
		policy.ErrorMsgPopup,
		policy.DisableWatermark,
		policy.ThemeMode))
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
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-3N9fs", "Errors.ResourceOwnerMissing")
	}
	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	existingPolicy := NewOrgLabelPolicyWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "Org-0K9dq", "Errors.Org.LabelPolicy.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		orgAgg,
		policy.PrimaryColor,
		policy.BackgroundColor,
		policy.WarnColor,
		policy.FontColor,
		policy.PrimaryColorDark,
		policy.BackgroundColorDark,
		policy.WarnColorDark,
		policy.FontColorDark,
		policy.HideLoginNameSuffix,
		policy.ErrorMsgPopup,
		policy.DisableWatermark,
		policy.ThemeMode)
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "Org-8nfSr", "Errors.Org.LabelPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
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
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-KKd4X", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-34mSE", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyActivatedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddLogoLabelPolicy(ctx context.Context, orgID string, upload *AssetUpload) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-KKd4X", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-23BMs", "Errors.Org.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-4N3nf", "Errors.Assets.Object.PutFailed")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyLogoAddedEvent(ctx, orgAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveLogoLabelPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-2FN8s", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-4MVsf", "Errors.Org.LabelPolicy.NotFound")
	}
	err = c.removeAsset(ctx, orgID, existingPolicy.LogoKey)
	if err != nil {
		return nil, err
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyLogoRemovedEvent(ctx, orgAgg, existingPolicy.LogoKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddIconLabelPolicy(ctx context.Context, orgID string, upload *AssetUpload) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-hMDs3", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-4nq2f", "Errors.Org.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-4BS7f", "Errors.Assets.Object.PutFailed")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyIconAddedEvent(ctx, orgAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIconLabelPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-1nd0d", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-1nd9f", "Errors.Org.LabelPolicy.NotFound")
	}

	err = c.removeAsset(ctx, orgID, existingPolicy.IconKey)
	if err != nil {
		return nil, err
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyIconRemovedEvent(ctx, orgAgg, existingPolicy.IconKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddLogoDarkLabelPolicy(ctx context.Context, orgID string, upload *AssetUpload) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-67Ms2", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-QSqcd", "Errors.Org.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-3S7fN", "Errors.Assets.Object.PutFailed")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyLogoDarkAddedEvent(ctx, orgAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveLogoDarkLabelPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-4NF0d", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-0peQw", "Errors.Org.LabelPolicy.NotFound")
	}
	err = c.removeAsset(ctx, orgID, existingPolicy.LogoDarkKey)
	if err != nil {
		return nil, err
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyLogoDarkRemovedEvent(ctx, orgAgg, existingPolicy.LogoDarkKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddIconDarkLabelPolicy(ctx context.Context, orgID string, upload *AssetUpload) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-tzBfs", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-4Nf8s", "Errors.Org.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-4B7cs", "Errors.Assets.Object.PutFailed")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyIconDarkAddedEvent(ctx, orgAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIconDarkLabelPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-Mv9ds", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-3NFos", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyIconDarkRemovedEvent(ctx, orgAgg, existingPolicy.IconDarkKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddFontLabelPolicy(ctx context.Context, orgID string, upload *AssetUpload) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-1Nf9s", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-2M9fs", "Errors.Org.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-2f9fw", "Errors.Assets.Object.PutFailed")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyFontAddedEvent(ctx, orgAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveFontLabelPolicy(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	if orgID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-2n0fW", "Errors.ResourceOwnerMissing")
	}
	existingPolicy, err := c.orgLabelPolicyWriteModelByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "ORG-4n9SD", "Errors.Org.LabelPolicy.NotFound")
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, org.NewLabelPolicyFontRemovedEvent(ctx, orgAgg, existingPolicy.FontKey))
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
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-Mf9sf", "Errors.ResourceOwnerMissing")
	}
	existingPolicy := NewOrgLabelPolicyWriteModel(orgID)
	removeEvent, err := c.removeLabelPolicy(ctx, existingPolicy)
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, removeEvent)
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
		return nil, zerrors.ThrowNotFound(nil, "Org-3M9df", "Errors.Org.LabelPolicy.NotFound")
	}

	err = c.removeAssetsFolder(ctx, existingPolicy.AggregateID, static.ObjectTypeStyling)
	if err != nil {
		return nil, err
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
	err = c.removeAssetsFolder(ctx, orgID, static.ObjectTypeStyling)
	if err != nil {
		return nil, err
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewLabelPolicyRemovedEvent(ctx, orgAgg), nil
}

func (c *Commands) removeLabelPolicyAssets(ctx context.Context, existingPolicy *OrgLabelPolicyWriteModel) (*org.LabelPolicyAssetsRemovedEvent, error) {
	err := c.removeAssetsFolder(ctx, existingPolicy.AggregateID, static.ObjectTypeStyling)
	if err != nil {
		return nil, err
	}
	orgAgg := OrgAggregateFromWriteModel(&existingPolicy.WriteModel)
	return org.NewLabelPolicyAssetsRemovedEvent(ctx, orgAgg), nil
}

func (c *Commands) orgLabelPolicyWriteModelByID(ctx context.Context, orgID string) (*OrgLabelPolicyWriteModel, error) {
	policy := NewOrgLabelPolicyWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

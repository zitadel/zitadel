package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ChangeDefaultLabelPolicy(ctx context.Context, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-0K9dq", "Errors.IAM.LabelPolicy.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		instanceAgg,
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
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-28fHe", "Errors.IAM.LabelPolicy.NotChanged")
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

func (c *Commands) ActivateDefaultLabelPolicy(ctx context.Context) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, prepareActivateDefaultLabelPolicy(instanceAgg))
	if err != nil {
		return nil, err
	}
	pushedEvents, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return pushedEventsToObjectDetails(pushedEvents), nil
}

func (c *Commands) AddLogoDefaultLabelPolicy(ctx context.Context, upload *AssetUpload) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-Qw0pd", "Errors.IAM.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INSTANCE-3m20c", "Errors.Assets.Object.PutFailed")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyLogoAddedEvent(ctx, instanceAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveLogoDefaultLabelPolicy(ctx context.Context) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-Xc8Kf", "Errors.IAM.LabelPolicy.NotFound")
	}

	err = c.removeAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.LogoKey)
	if err != nil {
		return nil, err
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyLogoRemovedEvent(ctx, instanceAgg, existingPolicy.LogoKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddIconDefaultLabelPolicy(ctx context.Context, upload *AssetUpload) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-1yMx0", "Errors.IAM.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INSTANCE-yxE4f", "Errors.Assets.Object.PutFailed")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyIconAddedEvent(ctx, instanceAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIconDefaultLabelPolicy(ctx context.Context) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-4M0qw", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.removeAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.IconKey)
	if err != nil {
		return nil, err
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyIconRemovedEvent(ctx, instanceAgg, existingPolicy.IconKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddLogoDarkDefaultLabelPolicy(ctx context.Context, upload *AssetUpload) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-ZR9fs", "Errors.Instance.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INSTANCE-4fMs9", "Errors.Assets.Object.PutFailed")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyLogoDarkAddedEvent(ctx, instanceAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveLogoDarkDefaultLabelPolicy(ctx context.Context) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-3FGds", "Errors.Instance.LabelPolicy.NotFound")
	}
	err = c.removeAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.LogoDarkKey)
	if err != nil {
		return nil, err
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyLogoDarkRemovedEvent(ctx, instanceAgg, existingPolicy.LogoDarkKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddIconDarkDefaultLabelPolicy(ctx context.Context, upload *AssetUpload) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-vMsf9", "Errors.Instance.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "INSTANCE-1cxM3", "Errors.Assets.Object.PutFailed")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyIconDarkAddedEvent(ctx, instanceAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveIconDarkDefaultLabelPolicy(ctx context.Context) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-2nc7F", "Errors.Instance.LabelPolicy.NotFound")
	}
	err = c.removeAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.IconDarkKey)
	if err != nil {
		return nil, err
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyIconDarkRemovedEvent(ctx, instanceAgg, existingPolicy.IconDarkKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddFontDefaultLabelPolicy(ctx context.Context, upload *AssetUpload) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-1N8fE", "Errors.Instance.LabelPolicy.NotFound")
	}
	asset, err := c.uploadAsset(ctx, upload)
	if err != nil {
		return nil, zerrors.ThrowInternal(nil, "INSTANCE-1N8fs", "Errors.Assets.Object.PutFailed")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyFontAddedEvent(ctx, instanceAgg, asset.Name))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) RemoveFontDefaultLabelPolicy(ctx context.Context) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-Tk0gw", "Errors.Instance.LabelPolicy.NotFound")
	}
	err = c.removeAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.FontKey)
	if err != nil {
		return nil, err
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyFontRemovedEvent(ctx, instanceAgg, existingPolicy.FontKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) defaultLabelPolicyWriteModelByID(ctx context.Context) (policy *InstanceLabelPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewInstanceLabelPolicyWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func (c *Commands) getDefaultLabelPolicy(ctx context.Context) (*domain.LabelPolicy, error) {
	policyWriteModel, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}
	policy := writeModelToLabelPolicy(&policyWriteModel.LabelPolicyWriteModel)
	policy.Default = true
	return policy, nil
}

// prepareAddDefaultLabelPolicy adds a default label policy, if none exists prior, and activates it directly
// this functions is only used on instance setup so the policy can be activated directly
func prepareAddDefaultLabelPolicy(
	a *instance.Aggregate,
	primaryColor,
	backgroundColor,
	warnColor,
	fontColor,
	primaryColorDark,
	backgroundColorDark,
	warnColorDark,
	fontColorDark string,
	hideLoginNameSuffix,
	errorMsgPopup,
	disableWatermark bool,
	themeMode domain.LabelPolicyThemeMode,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceLabelPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if writeModel.State == domain.PolicyStateActive {
				return nil, zerrors.ThrowAlreadyExists(nil, "INSTANCE-2B0ps", "Errors.Instance.LabelPolicy.AlreadyExists")
			}
			return []eventstore.Command{
				instance.NewLabelPolicyAddedEvent(ctx, &a.Aggregate,
					primaryColor,
					backgroundColor,
					warnColor,
					fontColor,
					primaryColorDark,
					backgroundColorDark,
					warnColorDark,
					fontColorDark,
					hideLoginNameSuffix,
					errorMsgPopup,
					disableWatermark,
					themeMode,
				),
				instance.NewLabelPolicyActivatedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}

func prepareActivateDefaultLabelPolicy(
	a *instance.Aggregate,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel := NewInstanceLabelPolicyWriteModel(ctx)
			events, err := filter(ctx, writeModel.Query())
			if err != nil {
				return nil, err
			}
			writeModel.AppendEvents(events...)
			if err = writeModel.Reduce(); err != nil {
				return nil, err
			}
			if !writeModel.State.Exists() {
				return nil, zerrors.ThrowNotFound(nil, "INSTANCE-6M23e", "Errors.Instance.LabelPolicy.NotFound")
			}
			return []eventstore.Command{
				instance.NewLabelPolicyActivatedEvent(ctx, &a.Aggregate),
			}, nil
		}, nil
	}
}

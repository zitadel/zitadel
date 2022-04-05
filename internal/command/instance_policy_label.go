package command

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultLabelPolicy(ctx context.Context, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	addedPolicy := NewInstanceLabelPolicyWriteModel(ctx)
	instanceAgg := InstanceAggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	event, err := c.addDefaultLabelPolicy(ctx, instanceAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLabelPolicy(&addedPolicy.LabelPolicyWriteModel), nil
}

func (c *Commands) addDefaultLabelPolicy(ctx context.Context, instanceAgg *eventstore.Aggregate, addedPolicy *InstanceLabelPolicyWriteModel, policy *domain.LabelPolicy) (eventstore.Command, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "INSTANCE-2B0ps", "Errors.IAM.LabelPolicy.AlreadyExists")
	}

	return instance.NewLabelPolicyAddedEvent(
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
		policy.DisableWatermark), nil

}

func (c *Commands) ChangeDefaultLabelPolicy(ctx context.Context, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	if err := policy.IsValid(); err != nil {
		return nil, err
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-0K9dq", "Errors.IAM.LabelPolicy.NotFound")
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
		policy.DisableWatermark)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "INSTANCE-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
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
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-6M23e", "Errors.IAM.LabelPolicy.NotFound")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyActivatedEvent(ctx, instanceAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) AddLogoDefaultLabelPolicy(ctx context.Context, storageKey string) (*domain.ObjectDetails, error) {
	if storageKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-3m20c", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-Qw0pd", "Errors.IAM.LabelPolicy.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyLogoAddedEvent(ctx, instanceAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-Xc8Kf", "Errors.IAM.LabelPolicy.NotFound")
	}

	err = c.RemoveAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.LogoKey)
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

func (c *Commands) AddIconDefaultLabelPolicy(ctx context.Context, storageKey string) (*domain.ObjectDetails, error) {
	if storageKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-yxE4f", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-1yMx0", "Errors.IAM.LabelPolicy.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyIconAddedEvent(ctx, instanceAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-4M0qw", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.IconKey)
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

func (c *Commands) AddLogoDarkDefaultLabelPolicy(ctx context.Context, storageKey string) (*domain.ObjectDetails, error) {
	if storageKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-4fMs9", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-ZR9fs", "Errors.IAM.LabelPolicy.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyLogoDarkAddedEvent(ctx, instanceAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-3FGds", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.LogoDarkKey)
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

func (c *Commands) AddIconDarkDefaultLabelPolicy(ctx context.Context, storageKey string) (*domain.ObjectDetails, error) {
	if storageKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-1cxM3", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-vMsf9", "Errors.IAM.LabelPolicy.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyIconDarkAddedEvent(ctx, instanceAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-2nc7F", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.IconDarkKey)
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

func (c *Commands) AddFontDefaultLabelPolicy(ctx context.Context, storageKey string) (*domain.ObjectDetails, error) {
	if storageKey == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "INSTANCE-1N8fs", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-1N8fE", "Errors.IAM.LabelPolicy.NotFound")
	}
	instanceAgg := InstanceAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.Push(ctx, instance.NewLabelPolicyFontAddedEvent(ctx, instanceAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "INSTANCE-Tk0gw", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, authz.GetInstance(ctx).InstanceID(), existingPolicy.FontKey)
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

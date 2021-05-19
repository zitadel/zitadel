package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) AddDefaultLabelPolicy(ctx context.Context, policy *domain.LabelPolicy) (*domain.LabelPolicy, error) {
	addedPolicy := NewIAMLabelPolicyWriteModel()
	iamAgg := IAMAggregateFromWriteModel(&addedPolicy.LabelPolicyWriteModel.WriteModel)
	event, err := c.addDefaultLabelPolicy(ctx, iamAgg, addedPolicy, policy)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, event)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addedPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToLabelPolicy(&addedPolicy.LabelPolicyWriteModel), nil
}

func (c *Commands) addDefaultLabelPolicy(ctx context.Context, iamAgg *eventstore.Aggregate, addedPolicy *IAMLabelPolicyWriteModel, policy *domain.LabelPolicy) (eventstore.EventPusher, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.State == domain.PolicyStateActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LabelPolicy.AlreadyExists")
	}

	return iam_repo.NewLabelPolicyAddedEvent(
		ctx,
		iamAgg,
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
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-0K9dq", "Errors.IAM.LabelPolicy.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	changedEvent, hasChanged := existingPolicy.NewChangedEvent(
		ctx,
		iamAgg,
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
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
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

func (c *Commands) ActivateDefaultLabelPolicy(ctx context.Context) (*domain.ObjectDetails, error) {
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-6M23e", "Errors.IAM.LabelPolicy.NotFound")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyActivatedEvent(ctx, iamAgg))
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-3m20c", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-Qw0pd", "Errors.IAM.LabelPolicy.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyLogoAddedEvent(ctx, iamAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "IAM-Xc8Kf", "Errors.IAM.LabelPolicy.NotFound")
	}

	err = c.RemoveAsset(ctx, domain.IAMID, existingPolicy.LogoKey)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyLogoRemovedEvent(ctx, iamAgg, existingPolicy.LogoKey))
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-yxE4f", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-1yMx0", "Errors.IAM.LabelPolicy.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyIconAddedEvent(ctx, iamAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "IAM-4M0qw", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, domain.IAMID, existingPolicy.IconKey)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyIconRemovedEvent(ctx, iamAgg, existingPolicy.IconKey))
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-4fMs9", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-ZR9fs", "Errors.IAM.LabelPolicy.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyLogoDarkAddedEvent(ctx, iamAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "IAM-3FGds", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, domain.IAMID, existingPolicy.LogoDarkKey)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyLogoDarkRemovedEvent(ctx, iamAgg, existingPolicy.LogoDarkKey))
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-1cxM3", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-vMsf9", "Errors.IAM.LabelPolicy.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyIconDarkAddedEvent(ctx, iamAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "IAM-2nc7F", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, domain.IAMID, existingPolicy.IconDarkKey)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyIconDarkRemovedEvent(ctx, iamAgg, existingPolicy.IconDarkKey))
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
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-1N8fs", "Errors.Assets.EmptyKey")
	}
	existingPolicy, err := c.defaultLabelPolicyWriteModelByID(ctx)
	if err != nil {
		return nil, err
	}

	if existingPolicy.State == domain.PolicyStateUnspecified || existingPolicy.State == domain.PolicyStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-1N8fE", "Errors.IAM.LabelPolicy.NotFound")
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyFontAddedEvent(ctx, iamAgg, storageKey))
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
		return nil, caos_errs.ThrowNotFound(nil, "IAM-Tk0gw", "Errors.IAM.LabelPolicy.NotFound")
	}
	err = c.RemoveAsset(ctx, domain.IAMID, existingPolicy.FontKey)
	if err != nil {
		return nil, err
	}
	iamAgg := IAMAggregateFromWriteModel(&existingPolicy.LabelPolicyWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, iam_repo.NewLabelPolicyFontRemovedEvent(ctx, iamAgg, existingPolicy.FontKey))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingPolicy, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingPolicy.LabelPolicyWriteModel.WriteModel), nil
}

func (c *Commands) defaultLabelPolicyWriteModelByID(ctx context.Context) (policy *IAMLabelPolicyWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := NewIAMLabelPolicyWriteModel()
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/instance"

	"github.com/caos/zitadel/internal/domain"
)

func (c *Commands) SetDefaultFeatures(ctx context.Context, features *domain.Features) (*domain.ObjectDetails, error) {
	existingFeatures := NewInstanceFeaturesWriteModel(ctx)
	setEvent, err := c.setDefaultFeatures(ctx, existingFeatures, features)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.Push(ctx, setEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFeatures, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFeatures.WriteModel), nil
}

func (c *Commands) setDefaultFeatures(ctx context.Context, existingFeatures *InstanceFeaturesWriteModel, features *domain.Features) (*instance.FeaturesSetEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	setEvent, hasChanged := existingFeatures.NewSetEvent(
		ctx,
		InstanceAggregateFromWriteModel(&existingFeatures.FeaturesWriteModel.WriteModel),
		features.TierName,
		features.TierDescription,
		features.State,
		features.StateDescription,
		features.AuditLogRetention,
		features.LoginPolicyFactors,
		features.LoginPolicyIDP,
		features.LoginPolicyPasswordless,
		features.LoginPolicyRegistration,
		features.LoginPolicyUsernameLogin,
		features.LoginPolicyPasswordReset,
		features.PasswordComplexityPolicy,
		features.LabelPolicyPrivateLabel,
		features.LabelPolicyWatermark,
		features.CustomDomain,
		features.PrivacyPolicy,
		features.MetadataUser,
		features.CustomTextMessage,
		features.CustomTextLogin,
		features.LockoutPolicy,
		features.ActionsAllowed,
		features.MaxActions,
	)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Features-GE4h2", "Errors.Features.NotChanged")
	}
	return setEvent, nil
}

func (c *Commands) getDefaultFeatures(ctx context.Context) (*domain.Features, error) {
	existingFeatures := NewInstanceFeaturesWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	features := writeModelToFeatures(&existingFeatures.FeaturesWriteModel)
	features.IsDefault = true
	return features, nil
}

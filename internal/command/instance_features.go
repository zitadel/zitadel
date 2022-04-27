package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func SetDefaultFeatures(
	a *instance.Aggregate,
	tierName,
	tierDescription string,
	state domain.FeaturesState,
	stateDescription string,
	retention time.Duration,
	loginPolicyFactors,
	loginPolicyIDP,
	loginPolicyPasswordless,
	loginPolicyRegistration,
	loginPolicyUsernameLogin,
	loginPolicyPasswordReset,
	passwordComplexityPolicy,
	labelPolicyPrivateLabel,
	labelPolicyWatermark,
	customDomain,
	privacyPolicy,
	metadataUser,
	customTextMessage,
	customTextLogin,
	lockoutPolicy bool,
	actionsAllowed domain.ActionsAllowed,
	maxActions int,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if !state.Valid() || state == domain.FeaturesStateUnspecified || state == domain.FeaturesStateRemoved {
			return nil, errors.ThrowInvalidArgument(nil, "INSTA-d3r1s", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			writeModel, err := defaultFeatures(ctx, filter)
			if err != nil {
				return nil, err
			}
			event, hasChanged := writeModel.NewSetEvent(ctx, &a.Aggregate,
				tierName,
				tierDescription,
				state,
				stateDescription,
				retention,
				loginPolicyFactors,
				loginPolicyIDP,
				loginPolicyPasswordless,
				loginPolicyRegistration,
				loginPolicyUsernameLogin,
				loginPolicyPasswordReset,
				passwordComplexityPolicy,
				labelPolicyPrivateLabel,
				labelPolicyWatermark,
				customDomain,
				privacyPolicy,
				metadataUser,
				customTextMessage,
				customTextLogin,
				lockoutPolicy,
				actionsAllowed,
				maxActions,
			)
			if !hasChanged {
				return nil, errors.ThrowPreconditionFailed(nil, "INSTA-GE4h2", "Errors.Features.NotChanged")
			}
			return []eventstore.Command{
				event,
			}, nil
		}, nil
	}
}

func defaultFeatures(ctx context.Context, filter preparation.FilterToQueryReducer) (*InstanceFeaturesWriteModel, error) {
	features := NewInstanceFeaturesWriteModel(ctx)
	events, err := filter(ctx, features.Query())
	if err != nil {
		return nil, err
	}
	if len(events) == 0 {
		return features, nil
	}
	features.AppendEvents(events...)
	err = features.Reduce()
	return features, err
}

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
		return nil, errors.ThrowPreconditionFailed(nil, "Features-GE4h2", "Errors.Features.NotChanged")
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

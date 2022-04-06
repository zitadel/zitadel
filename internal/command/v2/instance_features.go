package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
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

func defaultFeatures(ctx context.Context, filter preparation.FilterToQueryReducer) (*command.InstanceFeaturesWriteModel, error) {
	features := command.NewInstanceFeaturesWriteModel(ctx)
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

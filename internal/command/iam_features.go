package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/iam"

	"github.com/caos/zitadel/internal/domain"
)

func (c *Commands) SetDefaultFeatures(ctx context.Context, features *domain.Features) (*domain.ObjectDetails, error) {
	existingFeatures := NewIAMFeaturesWriteModel()
	setEvent, err := c.setDefaultFeatures(ctx, existingFeatures, features)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, setEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFeatures, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFeatures.WriteModel), nil
}

func (c *Commands) setDefaultFeatures(ctx context.Context, existingFeatures *IAMFeaturesWriteModel, features *domain.Features) (*iam.FeaturesSetEvent, error) {
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	setEvent, hasChanged := existingFeatures.NewSetEvent(
		ctx,
		IAMAggregateFromWriteModel(&existingFeatures.FeaturesWriteModel.WriteModel),
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
		features.PasswordComplexityPolicy,
		features.LabelPolicy,
		features.CustomDomain,
	)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Features-GE4h2", "Errors.Features.NotChanged")
	}
	return setEvent, nil
}

func (c *Commands) getDefaultFeatures(ctx context.Context) (*domain.Features, error) {
	existingFeatures := NewIAMFeaturesWriteModel()
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	return writeModelToFeatures(&existingFeatures.FeaturesWriteModel), nil
}

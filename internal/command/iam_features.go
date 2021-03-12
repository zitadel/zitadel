package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"

	"github.com/caos/zitadel/internal/domain"
)

func (c *Commands) SetDefaultFeatures(ctx context.Context, features *domain.Features) (*domain.ObjectDetails, error) {
	existingFeatures := NewIAMFeaturesWriteModel()
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	setEvent, hasChanged := existingFeatures.NewSetEvent(
		ctx,
		IAMAggregateFromWriteModel(&existingFeatures.FeaturesWriteModel.WriteModel),
		features.TierName,
		features.TierDescription,
		features.TierStatus,
		features.TierStatusDescription,
		features.LoginPolicyFactors,
		features.LoginPolicyIDP,
		features.LoginPolicyPasswordless,
		features.LoginPolicyRegistration,
		features.LoginPolicyUsernameLogin,
	)
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Features-GE4h2", "Errors.Features.NotChanged")
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

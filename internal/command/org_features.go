package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/repository/org"
)

func (c *Commands) SetOrgFeatures(ctx context.Context, resourceOwner string, features *domain.Features) (*domain.ObjectDetails, error) {
	existingFeatures := NewOrgFeaturesWriteModel(resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	setEvent, hasChanged := existingFeatures.NewSetEvent(
		ctx,
		OrgAggregateFromWriteModel(&existingFeatures.FeaturesWriteModel.WriteModel),
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

func (c *Commands) RemoveOrgFeatures(ctx context.Context, orgID string) (*domain.ObjectDetails, error) {
	existingFeatures := NewOrgFeaturesWriteModel(orgID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingFeatures)
	if err != nil {
		return nil, err
	}
	if existingFeatures.State == domain.FeaturesStateUnspecified || existingFeatures.State == domain.FeaturesStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "Features-Bg32G", "Errors.Features.NotFound")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingFeatures.FeaturesWriteModel.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, org.NewFeaturesRemovedEvent(ctx, orgAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingFeatures, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingFeatures.WriteModel), nil
}

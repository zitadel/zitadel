package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ChangeIDPJWTConfig(ctx context.Context, config *domain.JWTIDPConfig, resourceOwner string) (*domain.JWTIDPConfig, error) {
	if resourceOwner == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-ff8NF", "Errors.ResourceOwnerMissing")
	}
	if config.IDPConfigID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "Org-2n99f", "Errors.IDMissing")
	}
	existingConfig := NewOrgIDPJWTConfigWriteModel(config.IDPConfigID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "Org-67J9d", "Errors.Org.IDPConfig.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingConfig.WriteModel)
	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		orgAgg,
		config.IDPConfigID,
		config.JWTEndpoint,
		config.Issuer,
		config.KeysEndpoint,
		config.HeaderName)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "Org-2k9fs", "Errors.Org.IDPConfig.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToIDPJWTConfig(&existingConfig.JWTConfigWriteModel), nil
}

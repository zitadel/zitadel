package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

func (c *Commands) ChangeIDPJWTConfig(ctx context.Context, config *domain.JWTIDPConfig, resourceOwner string) (*domain.JWTIDPConfig, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-ff8NF", "Errors.ResourceOwnerMissing")
	}
	if config.IDPConfigID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-2n99f", "Errors.IDMissing")
	}
	existingConfig := NewOrgIDPJWTConfigWriteModel(config.IDPConfigID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "Org-67J9d", "Errors.Org.IDPConfig.AlreadyExists")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingConfig.WriteModel)
	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		orgAgg,
		config.IDPConfigID,
		config.JWTEndpoint,
		config.Issuer,
		config.KeysEndpoint)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-2k9fs", "Errors.Org.IDPConfig.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}

	return writeModelToIDPJWTConfig(&existingConfig.JWTConfigWriteModel), nil
}

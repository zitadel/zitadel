package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

func (c *Commands) ChangeDefaultIDPJWTConfig(ctx context.Context, config *domain.JWTIDPConfig) (*domain.JWTIDPConfig, error) {
	if config.IDPConfigID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-m9322", "Errors.IDMissing")
	}
	existingConfig := NewIAMIDPJWTConfigWriteModel(config.IDPConfigID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-2m00d", "Errors.IAM.IDPConfig.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingConfig.WriteModel)
	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		iamAgg,
		config.IDPConfigID,
		config.Issuer,
		config.KeysEndpoint)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-3n9gg", "Errors.IAM.IDPConfig.NotChanged")
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

package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

func (c *Commands) ChangeDefaultIDPOIDCConfig(ctx context.Context, config *domain.OIDCIDPConfig) (*domain.OIDCIDPConfig, error) {
	if config.IDPConfigID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "IAM-9djf8", "Errors.IDMissing")
	}
	existingConfig := NewIAMIDPOIDCConfigWriteModel(config.IDPConfigID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowNotFound(nil, "IAM-67J9d", "Errors.IAM.IDPConfig.AlreadyExists")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingConfig.WriteModel)
	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		iamAgg,
		config.IDPConfigID,
		config.ClientID,
		config.Issuer,
		config.ClientSecretString,
		c.idpConfigSecretCrypto,
		config.IDPDisplayNameMapping,
		config.UsernameMapping,
		config.Scopes...)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPOIDCConfig(&existingConfig.OIDCConfigWriteModel), nil
}

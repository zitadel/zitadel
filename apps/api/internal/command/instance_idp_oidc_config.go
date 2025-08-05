package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ChangeDefaultIDPOIDCConfig(ctx context.Context, config *domain.OIDCIDPConfig) (*domain.OIDCIDPConfig, error) {
	if config.IDPConfigID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTANCE-9djf8", "Errors.IDMissing")
	}
	existingConfig := NewInstanceIDPOIDCConfigWriteModel(ctx, config.IDPConfigID)
	err := c.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, zerrors.ThrowNotFound(nil, "INSTANCE-67J9d", "Errors.IAM.IDPConfig.AlreadyExists")
	}

	instanceAgg := InstanceAggregateFromWriteModel(&existingConfig.WriteModel)
	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		instanceAgg,
		config.IDPConfigID,
		config.ClientID,
		config.Issuer,
		config.AuthorizationEndpoint,
		config.TokenEndpoint,
		config.ClientSecretString,
		c.idpConfigEncryption,
		config.IDPDisplayNameMapping,
		config.UsernameMapping,
		config.Scopes...)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, zerrors.ThrowPreconditionFailed(nil, "INSTANCE-d8kwF", "Errors.IAM.IDPConfig.NotChanged")
	}

	pushedEvents, err := c.eventstore.Push(ctx, changedEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingConfig, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToIDPOIDCConfig(&existingConfig.OIDCConfigWriteModel), nil
}

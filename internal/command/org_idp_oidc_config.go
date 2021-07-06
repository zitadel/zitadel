package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
)

func (c *Commands) ChangeIDPOIDCConfig(ctx context.Context, config *domain.OIDCIDPConfig, resourceOwner string) (*domain.OIDCIDPConfig, error) {
	if resourceOwner == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-4n8f2", "Errors.ResourceOwnerMissing")
	}
	if config.IDPConfigID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "Org-66Qwj", "Errors.IDMissing")
	}
	existingConfig := NewOrgIDPOIDCConfigWriteModel(config.IDPConfigID, resourceOwner)
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
		config.ClientID,
		config.Issuer,
		config.AuthorizationEndpoint,
		config.TokenEndpoint,
		config.ClientSecretString,
		c.idpConfigSecretCrypto,
		config.IDPDisplayNameMapping,
		config.UsernameMapping,
		config.Scopes...)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9vs", "Errors.NoChanges")
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

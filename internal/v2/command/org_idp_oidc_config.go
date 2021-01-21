package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
)

func (r *CommandSide) ChangeIDPOIDCConfig(ctx context.Context, config *domain.OIDCIDPConfig) (*domain.OIDCIDPConfig, error) {
	existingConfig := NewOrgIDPOIDCConfigWriteModel(config.IDPConfigID, config.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowAlreadyExists(nil, "Org-67J9d", "Errors.Org.IDPConfig.AlreadyExists")
	}

	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		config.IDPConfigID,
		config.ClientID,
		config.Issuer,
		config.ClientSecretString,
		r.idpConfigSecretCrypto,
		config.IDPDisplayNameMapping,
		config.UsernameMapping,
		config.Scopes...)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "Org-4M9vs", "Errors.Org.LabelPolicy.NotChanged")
	}

	orgAgg := OrgAggregateFromWriteModel(&existingConfig.WriteModel)
	orgAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingConfig, orgAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToIDPOIDCConfig(&existingConfig.OIDCConfigWriteModel), nil
}

package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/business/domain"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (r *CommandSide) ChangeDefaultIDPOIDCConfig(ctx context.Context, config *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error) {
	existingConfig := NewIDPOIDCConfigWriteModel(config.AggregateID, config.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, existingConfig)
	if err != nil {
		return nil, err
	}

	if existingConfig.State == domain.IDPConfigStateRemoved || existingConfig.State == domain.IDPConfigStateUnspecified {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-67J9d", "Errors.IAM.IDPConfig.AlreadyExists")
	}

	changedEvent, hasChanged, err := existingConfig.NewChangedEvent(
		ctx,
		config.ClientID,
		config.Issuer,
		config.ClientSecretString,
		r.secretCrypto,
		domain.OIDCMappingField(config.IDPDisplayNameMapping),
		domain.OIDCMappingField(config.UsernameMapping),
		config.Scopes...)
	if err != nil {
		return nil, err
	}
	if !hasChanged {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-4M9vs", "Errors.IAM.LabelPolicy.NotChanged")
	}

	iamAgg := IAMAggregateFromWriteModel(&existingConfig.WriteModel)
	iamAgg.PushEvents(changedEvent)

	err = r.eventstore.PushAggregate(ctx, existingConfig, iamAgg)
	if err != nil {
		return nil, err
	}

	return writeModelToIDPOIDCConfig(&existingConfig.OIDCConfigWriteModel), nil
}

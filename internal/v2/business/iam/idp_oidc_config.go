package iam

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

func (r *Repository) ChangeIDPOIDCConfig(ctx context.Context, config *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error) {
	writeModel := iam.NewIDPOIDCConfigWriteModel(config.AggregateID, config.IDPConfigID)

	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	var clientSecret *crypto.CryptoValue

	if config.ClientSecretString != "" {
		clientSecret, err = crypto.Crypt([]byte(config.ClientSecretString), r.secretCrypto)
		if err != nil {
			return nil, err
		}
	}

	aggregate := iam.AggregateFromWriteModel(&writeModel.ConfigWriteModel.WriteModel).
		PushIDPOIDCConfigChanged(
			ctx,
			writeModel,
			config.ClientID,
			config.Issuer,
			clientSecret,
			oidc.MappingField(config.IDPDisplayNameMapping),
			oidc.MappingField(config.UsernameMapping),
			config.Scopes...)

	events, err := r.eventstore.PushAggregates(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	writeModel.AppendEvents(events...)
	if err = writeModel.Reduce(); err != nil {
		return nil, err
	}

	return writeModelToIDPOIDCConfig(&writeModel.ConfigWriteModel), nil
}

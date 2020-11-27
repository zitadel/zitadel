package iam

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idp"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

func (r *Repository) IDPConfigByID(ctx context.Context, iamID, idpConfigID string) (*iam_model.IDPConfigView, error) {
	idpConfig := iam.NewIDPConfigReadModel(iamID, idpConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpConfig)
	if err != nil {
		return nil, err
	}

	return readModelToIDPConfigView(idpConfig), nil
}

func (r *Repository) AddIDPConfig(ctx context.Context, config *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	if config.OIDCConfig == nil {
		return nil, errors.ThrowInvalidArgument(nil, "IAM-eUpQU", "Errors.idp.config.notset")
	}

	idpConfigID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	wm := iam.NewIDPConfigWriteModel(config.AggregateID, idpConfigID)
	err = r.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}

	clientSecret, err := crypto.Crypt([]byte(config.OIDCConfig.ClientSecretString), r.secretCrypto)
	if err != nil {
		return nil, err
	}

	aggregate := iam.AggregateFromWriteModel(&wm.WriteModel).
		PushIDPConfigAdded(ctx, idpConfigID, config.Name, idp.ConfigType(config.Type), idp.StylingType(config.StylingType)).
		PushIDPOIDCConfigAdded(
			ctx,
			config.OIDCConfig.ClientID,
			idpConfigID,
			config.OIDCConfig.Issuer,
			clientSecret,
			oidc.MappingField(config.OIDCConfig.IDPDisplayNameMapping),
			oidc.MappingField(config.OIDCConfig.UsernameMapping),
			config.OIDCConfig.Scopes...)

	events, err := r.eventstore.PushAggregates(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	if err = wm.AppendAndReduce(events...); err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(wm), nil
}

func (r *Repository) ChangeIDPConfig(ctx context.Context, config *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	writeModel := iam.NewIDPConfigWriteModel(config.AggregateID, config.IDPConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	aggregate := iam.AggregateFromWriteModel(&writeModel.WriteModel).
		PushIDPConfigChanged(
			ctx,
			writeModel,
			config.IDPConfigID,
			config.Name,
			idp.ConfigType(config.Type),
			idp.StylingType(config.StylingType))

	events, err := r.eventstore.PushAggregates(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	if err = writeModel.AppendAndReduce(events...); err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(writeModel), nil
}

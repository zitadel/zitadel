package command

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idp"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

func (r *CommandSide) AddDefaultIDPConfig(ctx context.Context, config *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	if config.OIDCConfig == nil {
		return nil, errors.ThrowInvalidArgument(nil, "IAM-eUpQU", "Errors.idp.config.notset")
	}

	idpConfigID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	//TODO: check name unique on aggregate

	clientSecret, err := crypto.Crypt([]byte(config.OIDCConfig.ClientSecretString), r.secretCrypto)
	if err != nil {
		return nil, err
	}

	writeModel, err := r.pushDefaultIDPWriteModel(ctx, config.AggregateID, idpConfigID, func(a *iam.Aggregate, _ *iam.IDPConfigWriteModel) *iam.Aggregate {
		return a.
			PushIDPConfigAdded(
				ctx,
				idpConfigID,
				config.Name,
				idp.ConfigType(config.Type),
				idp.StylingType(config.StylingType)).
			PushIDPOIDCConfigAdded(
				ctx,
				config.OIDCConfig.ClientID,
				idpConfigID,
				config.OIDCConfig.Issuer,
				clientSecret,
				oidc.MappingField(config.OIDCConfig.IDPDisplayNameMapping),
				oidc.MappingField(config.OIDCConfig.UsernameMapping),
				config.OIDCConfig.Scopes...)
	})
	if err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(writeModel), nil
}

func (r *CommandSide) ChangeDefaultIDPConfig(ctx context.Context, config *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	writeModel, err := r.pushDefaultIDPWriteModel(ctx, config.AggregateID, config.IDPConfigID, func(a *iam.Aggregate, writeModel *iam.IDPConfigWriteModel) *iam.Aggregate {
		return a.PushIDPConfigChanged(
			ctx,
			writeModel,
			config.IDPConfigID,
			config.Name,
			idp.ConfigType(config.Type),
			idp.StylingType(config.StylingType))
	})
	if err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(writeModel), nil
}

func (r *CommandSide) DeactivateDefaultIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	writeModel, err := r.pushDefaultIDPWriteModel(ctx, iamID, idpID, func(a *iam.Aggregate, _ *iam.IDPConfigWriteModel) *iam.Aggregate {
		return a.PushIDPConfigDeactivated(ctx, idpID)
	})
	if err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(writeModel), nil
}

func (r *CommandSide) ReactivateDefaultIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	writeModel, err := r.pushDefaultIDPWriteModel(ctx, iamID, idpID, func(a *iam.Aggregate, _ *iam.IDPConfigWriteModel) *iam.Aggregate {
		return a.PushIDPConfigReactivated(ctx, idpID)
	})
	if err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(writeModel), nil
}

func (r *CommandSide) RemoveDefaultIDPConfig(ctx context.Context, iamID, idpID string) (*iam_model.IDPConfig, error) {
	writeModel, err := r.pushDefaultIDPWriteModel(ctx, iamID, idpID, func(a *iam.Aggregate, _ *iam.IDPConfigWriteModel) *iam.Aggregate {
		return a.PushIDPConfigRemoved(ctx, idpID)
	})
	if err != nil {
		return nil, err
	}

	return writeModelToIDPConfig(writeModel), nil
}

func (r *CommandSide) pushDefaultIDPWriteModel(ctx context.Context, iamID, idpID string, eventSetter func(*iam.Aggregate, *iam.IDPConfigWriteModel) *iam.Aggregate) (*iam.IDPConfigWriteModel, error) {
	writeModel := iam.NewIDPConfigWriteModel(iamID, idpID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	aggregate := eventSetter(IAMAggregateFromWriteModel(&writeModel.WriteModel), writeModel)
	err = r.eventstore.PushAggregate(ctx, writeModel, aggregate)
	if err != nil {
		return nil, err
	}

	return writeModel, nil
}

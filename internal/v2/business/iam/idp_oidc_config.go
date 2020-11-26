package iam

import (
	"context"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idp"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
)

func (r *Repository) IDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	query := eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, iam.AggregateType).
		EventData(map[string]interface{}{
			"idpConfigId": idpConfigID,
		})

	idpConfig := new(iam.IDPConfigReadModel)

	events, err := r.eventstore.FilterEvents(ctx, query)
	if err != nil {
		return nil, err
	}

	idpConfig.AppendEvents(events...)
	if err = idpConfig.Reduce(); err != nil {
		return nil, err
	}

	return readModelToIDPConfigView(idpConfig), nil
}

func (r *Repository) AddIDPConfig(ctx context.Context, config *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	readModel, err := r.iamByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}

	idpConfigID, err := r.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	aggregate := iam.AggregateFromReadModel(readModel).
		PushIDPConfigAdded(ctx, idpConfigID, config.Name, idp.ConfigType(config.Type), idp.StylingType(config.StylingType))

	if config.OIDCConfig != nil {
		clientSecret, err := crypto.Crypt([]byte(config.OIDCConfig.ClientSecretString), r.secretCrypto)
		if err != nil {
			return nil, err
		}
		aggregate = aggregate.PushIDPOIDCConfigAdded(
			ctx,
			config.OIDCConfig.ClientID,
			idpConfigID,
			config.OIDCConfig.Issuer,
			clientSecret,
			oidc.MappingField(config.OIDCConfig.IDPDisplayNameMapping),
			oidc.MappingField(config.OIDCConfig.UsernameMapping),
			config.OIDCConfig.Scopes...)
	}

	events, err := r.eventstore.PushAggregates(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	if err = readModel.AppendAndReduce(events...); err != nil {
		return nil, err
	}

	idpConfig := readModel.IDPByID(idpConfigID)
	if idpConfig == nil {
		return nil, errors.ThrowInternal(nil, "IAM-stZYB", "Errors.Internal")
	}

	return readModelToIDPConfig(idpConfig), nil
}

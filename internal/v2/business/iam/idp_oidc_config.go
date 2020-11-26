package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
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
	iam, err := r.iamByID(ctx, config.AggregateID)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

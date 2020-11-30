package iam

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

func (r *Repository) AddIDPProviderToLoginPolicy(ctx context.Context, idpProvider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	writeModel := iam.NewLoginPolicyIDPProviderWriteModel(idpProvider.AggregateID, idpProvider.IdpConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	aggregate := iam.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicyIDPProviderAddedEvent(ctx, idpProvider.IdpConfigID, provider.Type(idpProvider.Type))

	if err = r.eventstore.PushAggregate(ctx, writeModel, aggregate); err != nil {
		return nil, err
	}

	return writeModelToIDPProvider(writeModel), nil
}

func (r *Repository) RemoveIDPProviderFromLoginPolicy(ctx context.Context, idpProvider *iam_model.IDPProvider) error {
	writeModel := iam.NewLoginPolicyIDPProviderWriteModel(idpProvider.AggregateID, idpProvider.IdpConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return err
	}
	aggregate := iam.AggregateFromWriteModel(&writeModel.WriteModel).
		PushLoginPolicyIDPProviderAddedEvent(ctx, idpProvider.IdpConfigID, provider.Type(idpProvider.Type))

	return r.eventstore.PushAggregate(ctx, writeModel, aggregate)
}

package query

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
)

func (r *Queries) DefaultIDPConfigByID(ctx context.Context, idpConfigID string) (domain.IDPConfig, error) {
	idpConfig := NewIAMIDPConfigReadModel(r.iamID, idpConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpConfig)
	if err != nil {
		return nil, err
	}

	return readModelToIDPConfigDomain(idpConfig), nil
}

package query

import (
	"context"
	"github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/repository/iam"
)

func (r *QuerySide) DefaultIDPConfigByID(ctx context.Context, iamID, idpConfigID string) (*model.IDPConfigView, error) {
	idpConfig := iam.NewIDPConfigReadModel(iamID, idpConfigID)
	err := r.eventstore.FilterToQueryReducer(ctx, idpConfig)
	if err != nil {
		return nil, err
	}

	return readModelToIDPConfigView(idpConfig), nil
}

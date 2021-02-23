package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/query"

	"github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository struct {
	IAMID string

	IAMV2QuerySide *query.QuerySide
}

func (repo *IAMRepository) GetIAM(ctx context.Context) (*model.IAM, error) {
	return repo.IAMV2QuerySide.IAMByID(ctx, repo.IAMID)
}

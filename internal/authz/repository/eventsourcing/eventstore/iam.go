package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/query"

	"github.com/caos/zitadel/internal/iam/model"
)

type IamRepo struct {
	IAMID string

	IAMV2Query *query.QuerySide
}

func (repo *IamRepo) Health(ctx context.Context) error {
	return nil
}

func (repo *IamRepo) IamByID(ctx context.Context) (*model.IAM, error) {
	return repo.IAMV2Query.IAMByID(ctx, repo.IAMID)
}

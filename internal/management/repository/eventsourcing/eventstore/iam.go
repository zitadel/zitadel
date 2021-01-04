package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/v2/query"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMRepository struct {
	IAMV2Query *query.QuerySide
}

func (repo *IAMRepository) IAMByID(ctx context.Context, id string) (*iam_model.IAM, error) {
	return repo.IAMV2Query.IAMByID(ctx, id)
}

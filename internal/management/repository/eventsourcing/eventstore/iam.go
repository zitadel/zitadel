package eventstore

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/query"
)

type IAMRepository struct {
	IAMV2Query *query.Queries
}

func (repo *IAMRepository) IAMByID(ctx context.Context, id string) (*iam_model.IAM, error) {
	return repo.IAMV2Query.IAMByID(ctx, id)
}

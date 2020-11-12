package eventstore

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_business "github.com/caos/zitadel/internal/v2/business/iam"
)

type IAMRepository struct {
	IAMV2 *iam_business.Repository
}

func (repo *IAMRepository) IAMByID(ctx context.Context, id string) (*iam_model.IAM, error) {
	return repo.IAMV2.IAMByID(ctx, id)
}

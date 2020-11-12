package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/iam/model"
	iam_business "github.com/caos/zitadel/internal/v2/business/iam"
)

type IAMRepository struct {
	IAMID string

	IAMV2 *iam_business.Repository
}

func (repo *IAMRepository) GetIAM(ctx context.Context) (*model.IAM, error) {
	return repo.IAMV2.IAMByID(ctx, repo.IAMID)
}

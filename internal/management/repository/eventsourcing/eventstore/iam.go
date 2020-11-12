package eventstore

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_business "github.com/caos/zitadel/internal/v2/business/iam"
)

type IAMRepository struct {
	IAMEvents *eventsourcing.IAMEventstore

	IAMV2 *iam_business.Repository
}

func (repo *IAMRepository) IAMByID(ctx context.Context, id string) (*iam_model.IAM, error) {
	if repo.IAMV2 != nil {
		return repo.IAMV2.IAMByID(ctx, id)
	}
	return repo.IAMEvents.IAMByID(ctx, id)
}

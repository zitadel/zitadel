package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/iam/model"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IAMRepository struct {
	IAMID     string
	IAMEvents *iam_event.IAMEventstore
}

func (repo *IAMRepository) GetIAM(ctx context.Context) (*model.IAM, error) {
	return repo.IAMEvents.IAMByID(ctx, repo.IAMID)
}

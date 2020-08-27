package eventstore

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IAMRepository struct {
	IAMEvents *eventsourcing.IAMEventstore
}

func (repo *IAMRepository) IAMByID(ctx context.Context, id string) (*iam_model.IAM, error) {
	return repo.IAMEvents.IAMByID(ctx, id)
}

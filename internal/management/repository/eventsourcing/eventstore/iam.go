package eventstore

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IamRepository struct {
	IamEvents *eventsourcing.IamEventstore
}

func (repo *IamRepository) IamByID(ctx context.Context, id string) (*iam_model.Iam, error) {
	return repo.IamEvents.IamByID(ctx, id)
}

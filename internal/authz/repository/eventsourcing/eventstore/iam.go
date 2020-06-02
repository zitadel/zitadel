package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/iam/model"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IamRepo struct {
	IamID     string
	IamEvents *iam_event.IamEventstore
}

func (repo *IamRepo) Health(ctx context.Context) error {
	return repo.IamEvents.Health(ctx)
}

func (repo *IamRepo) IamByID(ctx context.Context) (*model.Iam, error) {
	return repo.IamEvents.IamByID(ctx, repo.IamID)
}

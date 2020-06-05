package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/iam/model"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IamRepository struct {
	IamID     string
	IamEvents *iam_event.IamEventstore
}

func (repo *IamRepository) GetIam(ctx context.Context) (*model.Iam, error) {
	return repo.IamEvents.IamByID(ctx, repo.IamID)
}

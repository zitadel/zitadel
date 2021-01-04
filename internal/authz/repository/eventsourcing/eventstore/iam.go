package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/v2/query"

	"github.com/caos/zitadel/internal/iam/model"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
)

type IamRepo struct {
	IAMID     string
	IAMEvents *iam_event.IAMEventstore

	IAMV2Query *query.QuerySide
}

func (repo *IamRepo) Health(ctx context.Context) error {
	return repo.IAMEvents.Health(ctx)
}

func (repo *IamRepo) IamByID(ctx context.Context) (*model.IAM, error) {
	return repo.IAMV2Query.IAMByID(ctx, repo.IAMID)
}

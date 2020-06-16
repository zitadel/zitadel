package eventstore

import (
	"context"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/org/repository/view/model"
)

type OrgRepository struct {
	SearchLimit uint64
	*org_es.OrgEventstore
	View *auth_view.View
}

func (repo *OrgRepository) SearchOrgs(ctx context.Context, request *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	request.EnsureLimit(repo.SearchLimit)
	members, count, err := repo.View.SearchOrgs(request)
	if err != nil {
		return nil, err
	}
	return &org_model.OrgSearchResult{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.OrgsToModel(members),
	}, nil
}

package eventstore

import (
	"context"
	"github.com/caos/logging"
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
	sequence, err := repo.View.GetLatestOrgSequence()
	logging.Log("EVENT-7Udhz").OnError(err).Warn("could not read latest org sequence")
	members, count, err := repo.View.SearchOrgs(request)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgSearchResult{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.OrgsToModel(members),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

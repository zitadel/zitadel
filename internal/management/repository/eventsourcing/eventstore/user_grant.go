package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	grant_event "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
)

type UserGrantRepo struct {
	SearchLimit     uint64
	UserGrantEvents *grant_event.UserGrantEventStore
	View            *view.View
}

func (repo *UserGrantRepo) UserGrantByID(ctx context.Context, grantID string) (*grant_model.UserGrantView, error) {
	grant, err := repo.View.UserGrantByID(grantID)
	if err != nil {
		return nil, err
	}
	return model.UserGrantToModel(grant), nil
}

func (repo *UserGrantRepo) AddUserGrant(ctx context.Context, grant *grant_model.UserGrant) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.AddUserGrant(ctx, grant)
}

func (repo *UserGrantRepo) ChangeUserGrant(ctx context.Context, grant *grant_model.UserGrant) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.ChangeUserGrant(ctx, grant)
}

func (repo *UserGrantRepo) DeactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.DeactivateUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) ReactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	return repo.UserGrantEvents.ReactivateUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) RemoveUserGrant(ctx context.Context, grantID string) error {
	return repo.UserGrantEvents.RemoveUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) BulkAddUserGrant(ctx context.Context, grants ...*grant_model.UserGrant) error {
	return repo.UserGrantEvents.AddUserGrants(ctx, grants...)
}

func (repo *UserGrantRepo) BulkChangeUserGrant(ctx context.Context, grants ...*grant_model.UserGrant) error {
	return repo.UserGrantEvents.ChangeUserGrants(ctx, grants...)
}

func (repo *UserGrantRepo) BulkRemoveUserGrant(ctx context.Context, grantIDs ...string) error {
	return repo.UserGrantEvents.RemoveUserGrants(ctx, grantIDs...)
}

func (repo *UserGrantRepo) SearchUserGrants(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.UserGrantSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	grants, count, err := repo.View.SearchUserGrants(request)
	if err != nil {
		return nil, err
	}
	result := &grant_model.UserGrantSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.UserGrantsToModel(grants),
	}
	sequence, timestamp, err := repo.View.GetLatestUserGrantSequence()
	if err == nil {
		result.Sequence = sequence
		result.Timestamp = timestamp
	}
	return result, nil
}

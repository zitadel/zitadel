package eventstore

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	caos_errors "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	global_model "github.com/caos/zitadel/internal/model"
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
	err := checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
	if err != nil {
		return nil, err
	}
	return repo.UserGrantEvents.AddUserGrant(ctx, grant)
}

func (repo *UserGrantRepo) ChangeUserGrant(ctx context.Context, grant *grant_model.UserGrant) (*grant_model.UserGrant, error) {
	err := checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
	if err != nil {
		return nil, err
	}
	return repo.UserGrantEvents.ChangeUserGrant(ctx, grant)
}

func (repo *UserGrantRepo) DeactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	grant, err := repo.UserGrantByID(ctx, grantID)
	if err != nil {
		return nil, err
	}
	err = checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
	if err != nil {
		return nil, err
	}
	return repo.UserGrantEvents.DeactivateUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) ReactivateUserGrant(ctx context.Context, grantID string) (*grant_model.UserGrant, error) {
	grant, err := repo.UserGrantByID(ctx, grantID)
	if err != nil {
		return nil, err
	}
	err = checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
	if err != nil {
		return nil, err
	}
	return repo.UserGrantEvents.ReactivateUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) RemoveUserGrant(ctx context.Context, grantID string) error {
	grant, err := repo.UserGrantByID(ctx, grantID)
	if err != nil {
		return err
	}
	err = checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
	if err != nil {
		return err
	}
	return repo.UserGrantEvents.RemoveUserGrant(ctx, grantID)
}

func (repo *UserGrantRepo) BulkAddUserGrant(ctx context.Context, grants ...*grant_model.UserGrant) error {
	for _, grant := range grants {
		err := checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
		if err != nil {
			return err
		}
	}
	return repo.UserGrantEvents.AddUserGrants(ctx, grants...)
}

func (repo *UserGrantRepo) BulkChangeUserGrant(ctx context.Context, grants ...*grant_model.UserGrant) error {
	for _, grant := range grants {
		err := checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
		if err != nil {
			return err
		}
	}
	return repo.UserGrantEvents.ChangeUserGrants(ctx, grants...)
}

func (repo *UserGrantRepo) BulkRemoveUserGrant(ctx context.Context, grantIDs ...string) error {
	for _, grantID := range grantIDs {
		grant, err := repo.UserGrantByID(ctx, grantID)
		if err != nil {
			return err
		}
		err = checkExplicitPermission(ctx, grant.GrantID, grant.ProjectID)
		if err != nil {
			return err
		}
	}
	return repo.UserGrantEvents.RemoveUserGrants(ctx, grantIDs...)
}

func (repo *UserGrantRepo) SearchUserGrants(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.UserGrantSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestUserGrantSequence()
	logging.Log("EVENT-5Viwf").OnError(err).Warn("could not read latest user grant sequence")

	permissions := authz.GetPermissionsFromCtx(ctx)
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetPermissionCtxIDs(permissions)
		if _, q := request.GetSearchQuery(grant_model.UserGrantSearchKeyProjectID); q != nil {
			containsID := false
			for _, id := range ids {
				if id == q.Value {
					containsID = true
					break
				}
			}
			if !containsID {
				result := &grant_model.UserGrantSearchResponse{
					Offset:      request.Offset,
					Limit:       request.Limit,
					TotalResult: uint64(0),
					Result:      []*grant_model.UserGrantView{},
				}
				if err == nil {
					result.Sequence = sequence.CurrentSequence
					result.Timestamp = sequence.CurrentTimestamp
				}
				return result, nil
			}
		} else {
			request.Queries = append(request.Queries, &grant_model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyProjectID, Method: global_model.SearchMethodIsOneOf, Value: ids})
		}
	}

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
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.CurrentTimestamp
	}
	return result, nil
}

func checkExplicitPermission(ctx context.Context, grantID, projectID string) error {
	permissions := authz.GetPermissionsFromCtx(ctx)
	if !authz.HasGlobalPermission(permissions) {
		ids := authz.GetPermissionCtxIDs(permissions)
		containsID := false
		if grantID != "" {
			containsID = listContainsID(ids, grantID)
			if containsID {
				return nil
			}
		}
		containsID = listContainsID(ids, projectID)
		if !containsID {
			return caos_errors.ThrowPermissionDenied(nil, "EVENT-Shu7e", "Errors.UserGrant.NoPermissionForProject")
		}
	}
	return nil
}

func listContainsID(ids []string, id string) bool {
	containsID := false
	for _, i := range ids {
		if i == id {
			containsID = true
			break
		}
	}
	return containsID
}

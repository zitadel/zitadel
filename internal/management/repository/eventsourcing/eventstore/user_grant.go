package eventstore

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	global_model "github.com/caos/zitadel/internal/model"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	grant_event "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing"
	"github.com/caos/zitadel/internal/usergrant/repository/view/model"
	"github.com/caos/zitadel/internal/view/repository"
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

func (repo *UserGrantRepo) UserGrantsByProjectID(ctx context.Context, projectID string) ([]*grant_model.UserGrantView, error) {
	grants, err := repo.View.UserGrantsByProjectID(projectID)
	if err != nil {
		return nil, err
	}
	return model.UserGrantsToModel(grants), nil
}

func (repo *UserGrantRepo) UserGrantsByUserID(ctx context.Context, userID string) ([]*grant_model.UserGrantView, error) {
	grants, err := repo.View.UserGrantsByUserID(userID)
	if err != nil {
		return nil, err
	}
	return model.UserGrantsToModel(grants), nil
}

func (repo *UserGrantRepo) SearchUserGrants(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.UserGrantSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, sequenceErr := repo.View.GetLatestUserGrantSequence("")
	logging.Log("EVENT-5Viwf").OnError(sequenceErr).Warn("could not read latest user grant sequence")

	result := handleSearchUserGrantPermissions(ctx, request, sequence)
	if result != nil {
		return result, nil
	}

	grants, count, err := repo.View.SearchUserGrants(request)
	if err != nil {
		return nil, err
	}

	result = &grant_model.UserGrantSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.UserGrantsToModel(grants),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func handleSearchUserGrantPermissions(ctx context.Context, request *grant_model.UserGrantSearchRequest, sequence *repository.CurrentSequence) *grant_model.UserGrantSearchResponse {
	permissions := authz.GetAllPermissionsFromCtx(ctx)
	if authz.HasGlobalExplicitPermission(permissions, projectReadPerm) {
		return nil
	}

	ids := authz.GetExplicitPermissionCtxIDs(permissions, projectReadPerm)
	if _, query := request.GetSearchQuery(grant_model.UserGrantSearchKeyGrantID); query != nil {
		result := checkContainsPermID(ids, query, request, sequence)
		if result != nil {
			return result
		}
		return nil
	}
	if _, query := request.GetSearchQuery(grant_model.UserGrantSearchKeyProjectID); query != nil {
		result := checkContainsPermID(ids, query, request, sequence)
		if result != nil {
			return result
		}
	}
	request.Queries = append(request.Queries, &grant_model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyProjectID, Method: global_model.SearchMethodIsOneOf, Value: ids})
	return nil
}

func checkContainsPermID(ids []string, query *grant_model.UserGrantSearchQuery, request *grant_model.UserGrantSearchRequest, sequence *repository.CurrentSequence) *grant_model.UserGrantSearchResponse {
	containsID := false
	for _, id := range ids {
		if id == query.Value {
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
		if sequence != nil {
			result.Sequence = sequence.CurrentSequence
			result.Timestamp = sequence.LastSuccessfulSpoolerRun
		}
		return result
	}
	return nil
}

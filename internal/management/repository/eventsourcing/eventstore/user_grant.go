package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/domain"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/internal/view/repository"
)

type UserGrantRepo struct {
	SearchLimit     uint64
	View            *view.View
	PrefixAvatarURL string
}

func (repo *UserGrantRepo) SearchUserGrants(ctx context.Context, request *grant_model.UserGrantSearchRequest) (*grant_model.UserGrantSearchResponse, error) {
	// err := request.EnsureLimit(repo.SearchLimit)
	// if err != nil {
	// 	return nil, err
	// }
	// result := handleSearchUserGrantPermissions(ctx, request, sequence)
	// if result != nil {
	// 	return result, nil
	// }

	// grants, count, err := repo.View.SearchUserGrants(request)
	// if err != nil {
	// 	return nil, err
	// }

	// result = &grant_model.UserGrantSearchResponse{
	// 	Offset:      request.Offset,
	// 	Limit:       request.Limit,
	// 	TotalResult: count,
	// 	Result:      model.UserGrantsToModel(grants, repo.PrefixAvatarURL),
	// }
	// if sequenceErr == nil {
	// 	result.Sequence = sequence.CurrentSequence
	// 	result.Timestamp = sequence.LastSuccessfulSpoolerRun
	// }

	result := &grant_model.UserGrantSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: 0,
		Result:      nil,
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
	request.Queries = append(request.Queries, &grant_model.UserGrantSearchQuery{Key: grant_model.UserGrantSearchKeyProjectID, Method: domain.SearchMethodIsOneOf, Value: ids})
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
	if containsID {
		return nil
	}
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

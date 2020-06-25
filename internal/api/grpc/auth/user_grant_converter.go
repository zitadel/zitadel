package auth

import (
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/pkg/auth/grpc"
)

func userGrantSearchRequestsToModel(request *grpc.UserGrantSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: userGrantSearchQueriesToModel(request.Queries),
	}
}

func userGrantSearchQueriesToModel(queries []*grpc.UserGrantSearchQuery) []*grant_model.UserGrantSearchQuery {
	converted := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userGrantSearchQueryToModel(q)
	}
	return converted
}

func userGrantSearchQueryToModel(query *grpc.UserGrantSearchQuery) *grant_model.UserGrantSearchQuery {
	return &grant_model.UserGrantSearchQuery{
		Key:    userGrantSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userGrantSearchKeyToModel(key grpc.UserGrantSearchKey) grant_model.UserGrantSearchKey {
	switch key {
	case grpc.UserGrantSearchKey_UserGrantSearchKey_ORG_ID:
		return grant_model.UserGrantSearchKeyResourceOwner
	case grpc.UserGrantSearchKey_UserGrantSearchKey_PROJECT_ID:
		return grant_model.UserGrantSearchKeyProjectID
	default:
		return grant_model.UserGrantSearchKeyUnspecified
	}
}

func myProjectOrgSearchRequestRequestsToModel(request *grpc.MyProjectOrgSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:        request.Offset,
		Limit:         request.Limit,
		Asc:           request.Asc,
		SortingColumn: grant_model.UserGrantSearchKeyResourceOwner,
		Queries:       myProjectOrgSearchQueriesToModel(request.Queries),
	}
}

func myProjectOrgSearchQueriesToModel(queries []*grpc.MyProjectOrgSearchQuery) []*grant_model.UserGrantSearchQuery {
	converted := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = myProjectOrgSearchQueryToModel(q)
	}
	return converted
}

func myProjectOrgSearchQueryToModel(query *grpc.MyProjectOrgSearchQuery) *grant_model.UserGrantSearchQuery {
	return &grant_model.UserGrantSearchQuery{
		Key:    myProjectOrgSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func myProjectOrgSearchKeyToModel(key grpc.MyProjectOrgSearchKey) grant_model.UserGrantSearchKey {
	switch key {
	case grpc.MyProjectOrgSearchKey_MYPROJECTORGSEARCHKEY_ORG_NAME:
		return grant_model.UserGrantSearchKeyOrgName
	default:
		return grant_model.UserGrantSearchKeyUnspecified
	}
}

func userGrantSearchResponseFromModel(response *grant_model.UserGrantSearchResponse) *grpc.UserGrantSearchResponse {
	return &grpc.UserGrantSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      userGrantViewsFromModel(response.Result),
	}
}

func userGrantViewsFromModel(users []*grant_model.UserGrantView) []*grpc.UserGrantView {
	converted := make([]*grpc.UserGrantView, len(users))
	for i, user := range users {
		converted[i] = userGrantViewFromModel(user)
	}
	return converted
}

func userGrantViewFromModel(grant *grant_model.UserGrantView) *grpc.UserGrantView {
	return &grpc.UserGrantView{
		UserId:    grant.UserID,
		OrgId:     grant.ResourceOwner,
		OrgName:   grant.OrgName,
		ProjectId: grant.ProjectID,
		Roles:     grant.RoleKeys,
	}
}

func projectOrgSearchResponseFromModel(response *grant_model.ProjectOrgSearchResponse) *grpc.MyProjectOrgSearchResponse {
	return &grpc.MyProjectOrgSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectOrgsFromModel(response.Result),
	}
}

func projectOrgsFromModel(projectOrgs []*grant_model.Org) []*grpc.Org {
	converted := make([]*grpc.Org, len(projectOrgs))
	for i, org := range projectOrgs {
		converted[i] = projectOrgFromModel(org)
	}
	return converted
}

func projectOrgFromModel(org *grant_model.Org) *grpc.Org {
	return &grpc.Org{
		Id:   org.OrgID,
		Name: org.OrgName,
	}
}

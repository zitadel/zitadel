package grpc

import (
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
)

func userGrantSearchRequestsToModel(request *UserGrantSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: userGrantSearchQueriesToModel(request.Queries),
	}
}

func userGrantSearchQueriesToModel(queries []*UserGrantSearchQuery) []*grant_model.UserGrantSearchQuery {
	converted := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userGrantSearchQueryToModel(q)
	}
	return converted
}

func userGrantSearchQueryToModel(query *UserGrantSearchQuery) *grant_model.UserGrantSearchQuery {
	return &grant_model.UserGrantSearchQuery{
		Key:    userGrantSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userGrantSearchKeyToModel(key UserGrantSearchKey) grant_model.UserGrantSearchKey {
	switch key {
	case UserGrantSearchKey_UserGrantSearchKey_ORG_ID:
		return grant_model.USERGRANTSEARCHKEY_RESOURCEOWNER
	case UserGrantSearchKey_UserGrantSearchKey_PROJECT_ID:
		return grant_model.USERGRANTSEARCHKEY_PROJECT_ID
	case UserGrantSearchKey_UserGrantSearchKey_USER_ID:
		return grant_model.USERGRANTSEARCHKEY_USER_ID
	default:
		return grant_model.USERGRANTSEARCHKEY_UNSPECIFIED
	}
}

func myProjectOrgSearchRequestRequestsToModel(request *MyProjectOrgSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:        request.Offset,
		Limit:         request.Limit,
		Asc:           request.Asc,
		SortingColumn: grant_model.USERGRANTSEARCHKEY_RESOURCEOWNER,
		Queries:       myProjectOrgSearchQueriesToModel(request.Queries),
	}
}

func myProjectOrgSearchQueriesToModel(queries []*MyProjectOrgSearchQuery) []*grant_model.UserGrantSearchQuery {
	converted := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = myProjectOrgSearchQueryToModel(q)
	}
	return converted
}

func myProjectOrgSearchQueryToModel(query *MyProjectOrgSearchQuery) *grant_model.UserGrantSearchQuery {
	return &grant_model.UserGrantSearchQuery{
		Key:    myProjectOrgSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func myProjectOrgSearchKeyToModel(key MyProjectOrgSearchKey) grant_model.UserGrantSearchKey {
	switch key {
	case MyProjectOrgSearchKey_MYPROJECTORGSEARCHKEY_ORG_NAME:
		return grant_model.USERGRANTSEARCHKEY_ORG_NAME
	default:
		return grant_model.USERGRANTSEARCHKEY_UNSPECIFIED
	}
}

func userGrantSearchResponseFromModel(response *grant_model.UserGrantSearchResponse) *UserGrantSearchResponse {
	return &UserGrantSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      userGrantViewsFromModel(response.Result),
	}
}

func userGrantViewsFromModel(users []*grant_model.UserGrantView) []*UserGrantView {
	converted := make([]*UserGrantView, len(users))
	for i, user := range users {
		converted[i] = userGrantViewFromModel(user)
	}
	return converted
}

func userGrantViewFromModel(grant *grant_model.UserGrantView) *UserGrantView {
	return &UserGrantView{
		UserId:    grant.UserID,
		OrgId:     grant.ResourceOwner,
		OrgName:   grant.OrgName,
		ProjectId: grant.ProjectID,
		Roles:     grant.RoleKeys,
	}
}

func projectOrgSearchResponseFromModel(response *grant_model.ProjectOrgSearchResponse) *MyProjectOrgSearchResponse {
	return &MyProjectOrgSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectOrgsFromModel(response.Result),
	}
}

func projectOrgsFromModel(projectOrgs []*grant_model.Org) []*Org {
	converted := make([]*Org, len(projectOrgs))
	for i, org := range projectOrgs {
		converted[i] = projectOrgFromModel(org)
	}
	return converted
}

func projectOrgFromModel(org *grant_model.Org) *Org {
	return &Org{
		Id:   org.OrgID,
		Name: org.OrgName,
	}
}

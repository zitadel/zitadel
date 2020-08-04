package auth

import (
	"github.com/caos/logging"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/pkg/grpc/auth"
	"github.com/golang/protobuf/ptypes"
)

func userGrantSearchRequestsToModel(request *auth.UserGrantSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:  request.Offset,
		Limit:   request.Limit,
		Queries: userGrantSearchQueriesToModel(request.Queries),
	}
}

func userGrantSearchQueriesToModel(queries []*auth.UserGrantSearchQuery) []*grant_model.UserGrantSearchQuery {
	converted := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userGrantSearchQueryToModel(q)
	}
	return converted
}

func userGrantSearchQueryToModel(query *auth.UserGrantSearchQuery) *grant_model.UserGrantSearchQuery {
	return &grant_model.UserGrantSearchQuery{
		Key:    userGrantSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userGrantSearchKeyToModel(key auth.UserGrantSearchKey) grant_model.UserGrantSearchKey {
	switch key {
	case auth.UserGrantSearchKey_UserGrantSearchKey_ORG_ID:
		return grant_model.UserGrantSearchKeyResourceOwner
	case auth.UserGrantSearchKey_UserGrantSearchKey_PROJECT_ID:
		return grant_model.UserGrantSearchKeyProjectID
	default:
		return grant_model.UserGrantSearchKeyUnspecified
	}
}

func myProjectOrgSearchRequestRequestsToModel(request *auth.MyProjectOrgSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:        request.Offset,
		Limit:         request.Limit,
		Asc:           request.Asc,
		SortingColumn: grant_model.UserGrantSearchKeyResourceOwner,
		Queries:       myProjectOrgSearchQueriesToModel(request.Queries),
	}
}

func myProjectOrgSearchQueriesToModel(queries []*auth.MyProjectOrgSearchQuery) []*grant_model.UserGrantSearchQuery {
	converted := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = myProjectOrgSearchQueryToModel(q)
	}
	return converted
}

func myProjectOrgSearchQueryToModel(query *auth.MyProjectOrgSearchQuery) *grant_model.UserGrantSearchQuery {
	return &grant_model.UserGrantSearchQuery{
		Key:    myProjectOrgSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func myProjectOrgSearchKeyToModel(key auth.MyProjectOrgSearchKey) grant_model.UserGrantSearchKey {
	switch key {
	case auth.MyProjectOrgSearchKey_MYPROJECTORGSEARCHKEY_ORG_NAME:
		return grant_model.UserGrantSearchKeyOrgName
	default:
		return grant_model.UserGrantSearchKeyUnspecified
	}
}

func userGrantSearchResponseFromModel(response *grant_model.UserGrantSearchResponse) *auth.UserGrantSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-Lsp0d").OnError(err).Debug("unable to parse timestamp")

	return &auth.UserGrantSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            userGrantViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func userGrantViewsFromModel(users []*grant_model.UserGrantView) []*auth.UserGrantView {
	converted := make([]*auth.UserGrantView, len(users))
	for i, user := range users {
		converted[i] = userGrantViewFromModel(user)
	}
	return converted
}

func userGrantViewFromModel(grant *grant_model.UserGrantView) *auth.UserGrantView {
	return &auth.UserGrantView{
		UserId:    grant.UserID,
		OrgId:     grant.ResourceOwner,
		OrgName:   grant.OrgName,
		ProjectId: grant.ProjectID,
		Roles:     grant.RoleKeys,
		GrantId:   grant.GrantID,
	}
}

func projectOrgSearchResponseFromModel(response *grant_model.ProjectOrgSearchResponse) *auth.MyProjectOrgSearchResponse {
	return &auth.MyProjectOrgSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectOrgsFromModel(response.Result),
	}
}

func projectOrgsFromModel(projectOrgs []*grant_model.Org) []*auth.Org {
	converted := make([]*auth.Org, len(projectOrgs))
	for i, org := range projectOrgs {
		converted[i] = projectOrgFromModel(org)
	}
	return converted
}

func projectOrgFromModel(org *grant_model.Org) *auth.Org {
	return &auth.Org{
		Id:   org.OrgID,
		Name: org.OrgName,
	}
}

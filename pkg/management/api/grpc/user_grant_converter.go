package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/golang/protobuf/ptypes"
)

func usergrantFromModel(grant *grant_model.UserGrant) *UserGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-ki9ds").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-sl9ew").OnError(err).Debug("unable to parse timestamp")

	return &UserGrant{
		Id:           grant.AggregateID,
		UserId:       grant.UserID,
		State:        usergrantStateFromModel(grant.State),
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     grant.Sequence,
		ProjectId:    grant.ProjectID,
		RoleKeys:     grant.RoleKeys,
	}
}

func userGrantCreateToModel(u *UserGrantCreate) *grant_model.UserGrant {
	return &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		UserID:     u.UserId,
		ProjectID:  u.ProjectId,
		RoleKeys:   u.RoleKeys,
	}
}

func userGrantUpdateToModel(u *UserGrantUpdate) *grant_model.UserGrant {
	return &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.Id},
		RoleKeys:   u.RoleKeys,
	}
}

func projectUserGrantUpdateToModel(u *ProjectUserGrantUpdate) *grant_model.UserGrant {
	return &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.Id},
		RoleKeys:   u.RoleKeys,
	}
}

func projectGrantUserGrantCreateToModel(u *ProjectGrantUserGrantCreate) *grant_model.UserGrant {
	return &grant_model.UserGrant{
		UserID:    u.UserId,
		ProjectID: u.ProjectId,
		RoleKeys:  u.RoleKeys,
	}
}

func projectGrantUserGrantUpdateToModel(u *ProjectGrantUserGrantUpdate) *grant_model.UserGrant {
	return &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.Id},
		RoleKeys:   u.RoleKeys,
	}
}

func userGrantSearchRequestsToModel(project *UserGrantSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:  project.Offset,
		Limit:   project.Limit,
		Queries: userGrantSearchQueriesToModel(project.Queries),
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
	case UserGrantSearchKey_USERGRANTSEARCHKEY_ORG_ID:
		return grant_model.USERGRANTSEARCHKEY_RESOURCEOWNER
	case UserGrantSearchKey_USERGRANTSEARCHKEY_PROJECT_ID:
		return grant_model.USERGRANTSEARCHKEY_PROJECT_ID
	case UserGrantSearchKey_USERGRANTSEARCHKEY_USER_ID:
		return grant_model.USERGRANTSEARCHKEY_USER_ID
	case UserGrantSearchKey_USERGRANTSEARCHKEY_ROLE_KEY:
		return grant_model.USERGRANTSEARCHKEY_ROLE_KEY
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
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-dl9we").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-lpsg5").OnError(err).Debug("unable to parse timestamp")

	return &UserGrantView{
		Id:            grant.ID,
		State:         usergrantStateFromModel(grant.State),
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Sequence:      grant.Sequence,
		ResourceOwner: grant.ResourceOwner,
		UserName:      grant.UserName,
		FirstName:     grant.FirstName,
		LastName:      grant.LastName,
		Email:         grant.Email,
		ProjectName:   grant.ProjectName,
		OrgName:       grant.OrgName,
		OrgDomain:     grant.OrgDomain,
		RoleKeys:      grant.RoleKeys,
	}
}

func usergrantStateFromModel(state grant_model.UserGrantState) UserGrantState {
	switch state {
	case grant_model.USERGRANTSTATE_ACTIVE:
		return UserGrantState_USERGRANTSTATE_ACTIVE
	case grant_model.USERGRANTSTATE_INACTIVE:
		return UserGrantState_USERGRANTSTATE_INACTIVE
	default:
		return UserGrantState_USERGRANTSTATE_UNSPECIFIED
	}
}

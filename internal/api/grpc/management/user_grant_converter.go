package management

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/eventstore/models"
	grant_model "github.com/caos/zitadel/internal/usergrant/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func usergrantFromModel(grant *grant_model.UserGrant) *management.UserGrant {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-ki9ds").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-sl9ew").OnError(err).Debug("unable to parse timestamp")

	return &management.UserGrant{
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

func userGrantCreateToModel(u *management.UserGrantCreate) *grant_model.UserGrant {
	return &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.UserId},
		UserID:     u.UserId,
		ProjectID:  u.ProjectId,
		RoleKeys:   u.RoleKeys,
		GrantID:    u.GrantId,
	}
}

func userGrantUpdateToModel(u *management.UserGrantUpdate) *grant_model.UserGrant {
	return &grant_model.UserGrant{
		ObjectRoot: models.ObjectRoot{AggregateID: u.Id},
		RoleKeys:   u.RoleKeys,
	}
}

func userGrantRemoveBulkToModel(u *management.UserGrantRemoveBulk) []string {
	ids := make([]string, len(u.Ids))
	for i, id := range u.Ids {
		ids[i] = id
	}
	return ids
}

func userGrantSearchRequestsToModel(project *management.UserGrantSearchRequest) *grant_model.UserGrantSearchRequest {
	return &grant_model.UserGrantSearchRequest{
		Offset:  project.Offset,
		Limit:   project.Limit,
		Queries: userGrantSearchQueriesToModel(project.Queries),
	}
}

func userGrantSearchQueriesToModel(queries []*management.UserGrantSearchQuery) []*grant_model.UserGrantSearchQuery {
	converted := make([]*grant_model.UserGrantSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = userGrantSearchQueryToModel(q)
	}
	return converted
}

func userGrantSearchQueryToModel(query *management.UserGrantSearchQuery) *grant_model.UserGrantSearchQuery {
	return &grant_model.UserGrantSearchQuery{
		Key:    userGrantSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func userGrantSearchKeyToModel(key management.UserGrantSearchKey) grant_model.UserGrantSearchKey {
	switch key {
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_WITH_GRANTED:
		return grant_model.UserGrantSearchKeyWithGranted
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_PROJECT_ID:
		return grant_model.UserGrantSearchKeyProjectID
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_USER_ID:
		return grant_model.UserGrantSearchKeyUserID
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_ROLE_KEY:
		return grant_model.UserGrantSearchKeyRoleKey
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_GRANT_ID:
		return grant_model.UserGrantSearchKeyGrantID
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_USER_NAME:
		return grant_model.UserGrantSearchKeyUserName
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_FIRST_NAME:
		return grant_model.UserGrantSearchKeyFirstName
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_LAST_NAME:
		return grant_model.UserGrantSearchKeyLastName
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_EMAIL:
		return grant_model.UserGrantSearchKeyEmail
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_ORG_NAME:
		return grant_model.UserGrantSearchKeyOrgName
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_ORG_DOMAIN:
		return grant_model.UserGrantSearchKeyOrgDomain
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_PROJECT_NAME:
		return grant_model.UserGrantSearchKeyProjectName
	case management.UserGrantSearchKey_USERGRANTSEARCHKEY_DISPLAY_NAME:
		return grant_model.UserGrantSearchKeyDisplayName
	default:
		return grant_model.UserGrantSearchKeyUnspecified
	}
}

func userGrantSearchResponseFromModel(response *grant_model.UserGrantSearchResponse) *management.UserGrantSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-Wd7hs").OnError(err).Debug("unable to parse timestamp")
	return &management.UserGrantSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            userGrantViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func userGrantViewsFromModel(users []*grant_model.UserGrantView) []*management.UserGrantView {
	converted := make([]*management.UserGrantView, len(users))
	for i, user := range users {
		converted[i] = userGrantViewFromModel(user)
	}
	return converted
}

func userGrantViewFromModel(grant *grant_model.UserGrantView) *management.UserGrantView {
	creationDate, err := ptypes.TimestampProto(grant.CreationDate)
	logging.Log("GRPC-dl9we").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(grant.ChangeDate)
	logging.Log("GRPC-lpsg5").OnError(err).Debug("unable to parse timestamp")

	return &management.UserGrantView{
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
		OrgDomain:     grant.OrgPrimaryDomain,
		RoleKeys:      grant.RoleKeys,
		UserId:        grant.UserID,
		ProjectId:     grant.ProjectID,
		OrgId:         grant.ResourceOwner,
		DisplayName:   grant.DisplayName,
		GrantId:       grant.GrantID,
	}
}

func usergrantStateFromModel(state grant_model.UserGrantState) management.UserGrantState {
	switch state {
	case grant_model.UserGrantStateActive:
		return management.UserGrantState_USERGRANTSTATE_ACTIVE
	case grant_model.UserGrantStateInactive:
		return management.UserGrantState_USERGRANTSTATE_INACTIVE
	default:
		return management.UserGrantState_USERGRANTSTATE_UNSPECIFIED
	}
}

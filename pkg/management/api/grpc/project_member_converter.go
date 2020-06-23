package grpc

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/golang/protobuf/ptypes"
)

func projectMemberFromModel(member *proj_model.ProjectMember) *ProjectMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-kd8re").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-dlei3").OnError(err).Debug("unable to parse timestamp")

	return &ProjectMember{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		Roles:        member.Roles,
	}
}

func projectMemberAddToModel(member *ProjectMemberAdd) *proj_model.ProjectMember {
	return &proj_model.ProjectMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.Id,
		},
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}

func projectMemberChangeToModel(member *ProjectMemberChange) *proj_model.ProjectMember {
	return &proj_model.ProjectMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.Id,
		},
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}

func projectMemberSearchRequestsToModel(role *ProjectMemberSearchRequest) *proj_model.ProjectMemberSearchRequest {
	return &proj_model.ProjectMemberSearchRequest{
		Offset:  role.Offset,
		Limit:   role.Limit,
		Queries: projectMemberSearchQueriesToModel(role.Queries),
	}
}

func projectMemberSearchQueriesToModel(queries []*ProjectMemberSearchQuery) []*proj_model.ProjectMemberSearchQuery {
	converted := make([]*proj_model.ProjectMemberSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = projectMemberSearchQueryToModel(q)
	}
	return converted
}

func projectMemberSearchQueryToModel(query *ProjectMemberSearchQuery) *proj_model.ProjectMemberSearchQuery {
	return &proj_model.ProjectMemberSearchQuery{
		Key:    projectMemberSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectMemberSearchKeyToModel(key ProjectMemberSearchKey) proj_model.ProjectMemberSearchKey {
	switch key {
	case ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_EMAIL:
		return proj_model.ProjectMemberSearchKeyEmail
	case ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_FIRST_NAME:
		return proj_model.ProjectMemberSearchKeyFirstName
	case ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_LAST_NAME:
		return proj_model.ProjectMemberSearchKeyLastName
	case ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_USER_NAME:
		return proj_model.ProjectMemberSearchKeyUserName
	default:
		return proj_model.ProjectMemberSearchKeyUnspecified
	}
}

func projectMemberSearchResponseFromModel(response *proj_model.ProjectMemberSearchResponse) *ProjectMemberSearchResponse {
	return &ProjectMemberSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectMemberViewsFromModel(response.Result),
	}
}

func projectMemberViewsFromModel(members []*proj_model.ProjectMemberView) []*ProjectMemberView {
	converted := make([]*ProjectMemberView, len(members))
	for i, member := range members {
		converted[i] = projectMemberViewFromModel(member)
	}
	return converted
}

func projectMemberViewFromModel(member *proj_model.ProjectMemberView) *ProjectMemberView {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-sl9cs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-8iw2d").OnError(err).Debug("unable to parse timestamp")

	return &ProjectMemberView{
		UserId:       member.UserID,
		UserName:     member.UserName,
		Email:        member.Email,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		Roles:        member.Roles,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
	}
}

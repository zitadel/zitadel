package management

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/management/grpc"
)

func projectGrantMemberFromModel(member *proj_model.ProjectGrantMember) *grpc.ProjectGrantMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-7du3s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-8duew").OnError(err).Debug("unable to parse timestamp")

	return &grpc.ProjectGrantMember{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		Roles:        member.Roles,
	}
}

func projectGrantMemberAddToModel(member *grpc.ProjectGrantMemberAdd) *proj_model.ProjectGrantMember {
	return &proj_model.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.ProjectId,
		},
		GrantID: member.GrantId,
		UserID:  member.UserId,
		Roles:   member.Roles,
	}
}

func projectGrantMemberChangeToModel(member *grpc.ProjectGrantMemberChange) *proj_model.ProjectGrantMember {
	return &proj_model.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.ProjectId,
		},
		GrantID: member.GrantId,
		UserID:  member.UserId,
		Roles:   member.Roles,
	}
}

func projectGrantMemberSearchRequestsToModel(role *grpc.ProjectGrantMemberSearchRequest) *proj_model.ProjectGrantMemberSearchRequest {
	return &proj_model.ProjectGrantMemberSearchRequest{
		Offset:  role.Offset,
		Limit:   role.Limit,
		Queries: projectGrantMemberSearchQueriesToModel(role.Queries),
	}
}

func projectGrantMemberSearchQueriesToModel(queries []*grpc.ProjectGrantMemberSearchQuery) []*proj_model.ProjectGrantMemberSearchQuery {
	converted := make([]*proj_model.ProjectGrantMemberSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = projectGrantMemberSearchQueryToModel(q)
	}
	return converted
}

func projectGrantMemberSearchQueryToModel(query *grpc.ProjectGrantMemberSearchQuery) *proj_model.ProjectGrantMemberSearchQuery {
	return &proj_model.ProjectGrantMemberSearchQuery{
		Key:    projectGrantMemberSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectGrantMemberSearchKeyToModel(key grpc.ProjectGrantMemberSearchKey) proj_model.ProjectGrantMemberSearchKey {
	switch key {
	case grpc.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_EMAIL:
		return proj_model.ProjectGrantMemberSearchKeyEmail
	case grpc.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_FIRST_NAME:
		return proj_model.ProjectGrantMemberSearchKeyFirstName
	case grpc.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_LAST_NAME:
		return proj_model.ProjectGrantMemberSearchKeyLastName
	case grpc.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_USER_NAME:
		return proj_model.ProjectGrantMemberSearchKeyUserName
	default:
		return proj_model.ProjectGrantMemberSearchKeyUnspecified
	}
}

func projectGrantMemberSearchResponseFromModel(response *proj_model.ProjectGrantMemberSearchResponse) *grpc.ProjectGrantMemberSearchResponse {
	return &grpc.ProjectGrantMemberSearchResponse{
		Offset:      response.Offset,
		Limit:       response.Limit,
		TotalResult: response.TotalResult,
		Result:      projectGrantMemberViewsFromModel(response.Result),
	}
}

func projectGrantMemberViewsFromModel(roles []*proj_model.ProjectGrantMemberView) []*grpc.ProjectGrantMemberView {
	converted := make([]*grpc.ProjectGrantMemberView, len(roles))
	for i, role := range roles {
		converted[i] = projectGrantMemberViewFromModel(role)
	}
	return converted
}

func projectGrantMemberViewFromModel(member *proj_model.ProjectGrantMemberView) *grpc.ProjectGrantMemberView {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-los93").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-ski4e").OnError(err).Debug("unable to parse timestamp")

	return &grpc.ProjectGrantMemberView{
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

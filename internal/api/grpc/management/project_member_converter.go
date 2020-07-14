package management

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/eventstore/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func projectMemberFromModel(member *proj_model.ProjectMember) *management.ProjectMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-kd8re").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-dlei3").OnError(err).Debug("unable to parse timestamp")

	return &management.ProjectMember{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		Roles:        member.Roles,
	}
}

func projectMemberAddToModel(member *management.ProjectMemberAdd) *proj_model.ProjectMember {
	return &proj_model.ProjectMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.Id,
		},
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}

func projectMemberChangeToModel(member *management.ProjectMemberChange) *proj_model.ProjectMember {
	return &proj_model.ProjectMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.Id,
		},
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}

func projectMemberSearchRequestsToModel(member *management.ProjectMemberSearchRequest) *proj_model.ProjectMemberSearchRequest {
	return &proj_model.ProjectMemberSearchRequest{
		Offset:  member.Offset,
		Limit:   member.Limit,
		Queries: projectMemberSearchQueriesToModel(member.Queries),
	}
}

func projectMemberSearchQueriesToModel(queries []*management.ProjectMemberSearchQuery) []*proj_model.ProjectMemberSearchQuery {
	converted := make([]*proj_model.ProjectMemberSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = projectMemberSearchQueryToModel(q)
	}
	return converted
}

func projectMemberSearchQueryToModel(query *management.ProjectMemberSearchQuery) *proj_model.ProjectMemberSearchQuery {
	return &proj_model.ProjectMemberSearchQuery{
		Key:    projectMemberSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectMemberSearchKeyToModel(key management.ProjectMemberSearchKey) proj_model.ProjectMemberSearchKey {
	switch key {
	case management.ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_EMAIL:
		return proj_model.ProjectMemberSearchKeyEmail
	case management.ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_FIRST_NAME:
		return proj_model.ProjectMemberSearchKeyFirstName
	case management.ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_LAST_NAME:
		return proj_model.ProjectMemberSearchKeyLastName
	case management.ProjectMemberSearchKey_PROJECTMEMBERSEARCHKEY_USER_NAME:
		return proj_model.ProjectMemberSearchKeyUserName
	default:
		return proj_model.ProjectMemberSearchKeyUnspecified
	}
}

func projectMemberSearchResponseFromModel(response *proj_model.ProjectMemberSearchResponse) *management.ProjectMemberSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-LSo9j").OnError(err).Debug("unable to parse timestamp")
	return &management.ProjectMemberSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            projectMemberViewsFromModel(response.Result),
		ViewTimestamp:     timestamp,
		ProcessedSequence: response.Sequence,
	}
}

func projectMemberViewsFromModel(members []*proj_model.ProjectMemberView) []*management.ProjectMemberView {
	converted := make([]*management.ProjectMemberView, len(members))
	for i, member := range members {
		converted[i] = projectMemberViewFromModel(member)
	}
	return converted
}

func projectMemberViewFromModel(member *proj_model.ProjectMemberView) *management.ProjectMemberView {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-sl9cs").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-8iw2d").OnError(err).Debug("unable to parse timestamp")

	return &management.ProjectMemberView{
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

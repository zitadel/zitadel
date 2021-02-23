package management

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/model"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/caos/zitadel/internal/eventstore/v1/models"
	proj_model "github.com/caos/zitadel/internal/project/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func projectGrantMemberFromDomain(member *domain.ProjectGrantMember) *management.ProjectGrantMember {
	return &management.ProjectGrantMember{
		CreationDate: timestamppb.New(member.CreationDate),
		ChangeDate:   timestamppb.New(member.ChangeDate),
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		Roles:        member.Roles,
	}
}

func projectGrantMemberFromModel(member *proj_model.ProjectGrantMember) *management.ProjectGrantMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-7du3s").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-8duew").OnError(err).Debug("unable to parse timestamp")

	return &management.ProjectGrantMember{
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		Roles:        member.Roles,
	}
}

func projectGrantMemberAddToDomain(member *management.ProjectGrantMemberAdd) *domain.ProjectGrantMember {
	return &domain.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.ProjectId,
		},
		GrantID: member.GrantId,
		UserID:  member.UserId,
		Roles:   member.Roles,
	}
}

func projectGrantMemberChangeToDomain(member *management.ProjectGrantMemberChange) *domain.ProjectGrantMember {
	return &domain.ProjectGrantMember{
		ObjectRoot: models.ObjectRoot{
			AggregateID: member.ProjectId,
		},
		GrantID: member.GrantId,
		UserID:  member.UserId,
		Roles:   member.Roles,
	}
}

func projectGrantMemberSearchRequestsToModel(memberSearch *management.ProjectGrantMemberSearchRequest) *proj_model.ProjectGrantMemberSearchRequest {
	request := &proj_model.ProjectGrantMemberSearchRequest{
		Offset:  memberSearch.Offset,
		Limit:   memberSearch.Limit,
		Queries: projectGrantMemberSearchQueriesToModel(memberSearch.Queries),
	}
	request.Queries = append(request.Queries, &proj_model.ProjectGrantMemberSearchQuery{Key: proj_model.ProjectGrantMemberSearchKeyProjectID, Method: model.SearchMethodEquals, Value: memberSearch.ProjectId})
	request.Queries = append(request.Queries, &proj_model.ProjectGrantMemberSearchQuery{Key: proj_model.ProjectGrantMemberSearchKeyGrantID, Method: model.SearchMethodEquals, Value: memberSearch.GrantId})
	return request
}

func projectGrantMemberSearchQueriesToModel(queries []*management.ProjectGrantMemberSearchQuery) []*proj_model.ProjectGrantMemberSearchQuery {
	converted := make([]*proj_model.ProjectGrantMemberSearchQuery, len(queries))
	for i, q := range queries {
		converted[i] = projectGrantMemberSearchQueryToModel(q)
	}
	return converted
}

func projectGrantMemberSearchQueryToModel(query *management.ProjectGrantMemberSearchQuery) *proj_model.ProjectGrantMemberSearchQuery {
	return &proj_model.ProjectGrantMemberSearchQuery{
		Key:    projectGrantMemberSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func projectGrantMemberSearchKeyToModel(key management.ProjectGrantMemberSearchKey) proj_model.ProjectGrantMemberSearchKey {
	switch key {
	case management.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_EMAIL:
		return proj_model.ProjectGrantMemberSearchKeyEmail
	case management.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_FIRST_NAME:
		return proj_model.ProjectGrantMemberSearchKeyFirstName
	case management.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_LAST_NAME:
		return proj_model.ProjectGrantMemberSearchKeyLastName
	case management.ProjectGrantMemberSearchKey_PROJECTGRANTMEMBERSEARCHKEY_USER_NAME:
		return proj_model.ProjectGrantMemberSearchKeyUserName
	default:
		return proj_model.ProjectGrantMemberSearchKeyUnspecified
	}
}

func projectGrantMemberSearchResponseFromModel(response *proj_model.ProjectGrantMemberSearchResponse) *management.ProjectGrantMemberSearchResponse {
	timestamp, err := ptypes.TimestampProto(response.Timestamp)
	logging.Log("GRPC-MSn6g").OnError(err).Debug("unable to parse timestamp")
	return &management.ProjectGrantMemberSearchResponse{
		Offset:            response.Offset,
		Limit:             response.Limit,
		TotalResult:       response.TotalResult,
		Result:            projectGrantMemberViewsFromModel(response.Result),
		ProcessedSequence: response.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func projectGrantMemberViewsFromModel(roles []*proj_model.ProjectGrantMemberView) []*management.ProjectGrantMemberView {
	converted := make([]*management.ProjectGrantMemberView, len(roles))
	for i, role := range roles {
		converted[i] = projectGrantMemberViewFromModel(role)
	}
	return converted
}

func projectGrantMemberViewFromModel(member *proj_model.ProjectGrantMemberView) *management.ProjectGrantMemberView {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-los93").OnError(err).Debug("unable to parse timestamp")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-ski4e").OnError(err).Debug("unable to parse timestamp")

	return &management.ProjectGrantMemberView{
		UserId:       member.UserID,
		UserName:     member.UserName,
		Email:        member.Email,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		DisplayName:  member.DisplayName,
		Roles:        member.Roles,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     member.Sequence,
	}
}

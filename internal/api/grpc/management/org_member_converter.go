package management

import (
	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func addOrgMemberToModel(member *management.AddOrgMemberRequest) *org_model.OrgMember {
	memberModel := &org_model.OrgMember{
		UserID: member.UserId,
	}
	memberModel.Roles = member.Roles

	return memberModel
}

func changeOrgMemberToModel(member *management.ChangeOrgMemberRequest) *org_model.OrgMember {
	memberModel := &org_model.OrgMember{
		UserID: member.UserId,
	}
	memberModel.Roles = member.Roles

	return memberModel
}

func orgMemberFromModel(member *org_model.OrgMember) *management.OrgMember {
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-jC5wY").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-Nc2jJ").OnError(err).Debug("date parse failed")

	return &management.OrgMember{
		UserId:       member.UserID,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Roles:        member.Roles,
		Sequence:     member.Sequence,
	}
}

func orgMemberSearchRequestToModel(request *management.OrgMemberSearchRequest) *org_model.OrgMemberSearchRequest {
	return &org_model.OrgMemberSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: orgMemberSearchQueriesToModel(request.Queries),
	}
}

func orgMemberSearchQueriesToModel(queries []*management.OrgMemberSearchQuery) []*org_model.OrgMemberSearchQuery {
	modelQueries := make([]*org_model.OrgMemberSearchQuery, len(queries)+1)

	for i, query := range queries {
		modelQueries[i] = orgMemberSearchQueryToModel(query)
	}

	return modelQueries
}

func orgMemberSearchQueryToModel(query *management.OrgMemberSearchQuery) *org_model.OrgMemberSearchQuery {
	return &org_model.OrgMemberSearchQuery{
		Key:    orgMemberSearchKeyToModel(query.Key),
		Method: orgMemberSearchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func orgMemberSearchKeyToModel(key management.OrgMemberSearchKey) org_model.OrgMemberSearchKey {
	switch key {
	case management.OrgMemberSearchKey_ORGMEMBERSEARCHKEY_EMAIL:
		return org_model.OrgMemberSearchKeyEmail
	case management.OrgMemberSearchKey_ORGMEMBERSEARCHKEY_FIRST_NAME:
		return org_model.OrgMemberSearchKeyFirstName
	case management.OrgMemberSearchKey_ORGMEMBERSEARCHKEY_LAST_NAME:
		return org_model.OrgMemberSearchKeyLastName
	case management.OrgMemberSearchKey_ORGMEMBERSEARCHKEY_USER_ID:
		return org_model.OrgMemberSearchKeyUserID
	default:
		return org_model.OrgMemberSearchKeyUnspecified
	}
}

func orgMemberSearchMethodToModel(key management.SearchMethod) model.SearchMethod {
	switch key {
	case management.SearchMethod_SEARCHMETHOD_CONTAINS:
		return model.SearchMethodContains
	case management.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return model.SearchMethodContainsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_EQUALS:
		return model.SearchMethodEquals
	case management.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return model.SearchMethodEqualsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return model.SearchMethodStartsWith
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return model.SearchMethodStartsWithIgnoreCase
	default:
		return -1
	}
}

func orgMemberSearchResponseFromModel(resp *org_model.OrgMemberSearchResponse) *management.OrgMemberSearchResponse {
	timestamp, err := ptypes.TimestampProto(resp.Timestamp)
	logging.Log("GRPC-Swmr6").OnError(err).Debug("date parse failed")
	return &management.OrgMemberSearchResponse{
		Limit:             resp.Limit,
		Offset:            resp.Offset,
		TotalResult:       resp.TotalResult,
		Result:            orgMembersFromView(resp.Result),
		ProcessedSequence: resp.Sequence,
		ViewTimestamp:     timestamp,
	}
}
func orgMembersFromView(viewMembers []*org_model.OrgMemberView) []*management.OrgMemberView {
	members := make([]*management.OrgMemberView, len(viewMembers))

	for i, member := range viewMembers {
		members[i] = orgMemberFromView(member)
	}

	return members
}

func orgMemberFromView(member *org_model.OrgMemberView) *management.OrgMemberView {
	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-S9LAZ").OnError(err).Debug("unable to parse changedate")
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-oJN56").OnError(err).Debug("unable to parse creation date")

	return &management.OrgMemberView{
		ChangeDate:   changeDate,
		CreationDate: creationDate,
		Roles:        member.Roles,
		Sequence:     member.Sequence,
		UserId:       member.UserID,
		UserName:     member.UserName,
		Email:        member.Email,
		FirstName:    member.FirstName,
		LastName:     member.LastName,
		DisplayName:  member.DisplayName,
	}
}

package management

import (
	"context"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func addOrgMemberToDomain(ctx context.Context, member *management.AddOrgMemberRequest) *domain.Member {
	return domain.NewMember(authz.GetCtxData(ctx).OrgID, member.UserId, member.Roles...)
}

func changeOrgMemberToModel(ctx context.Context, member *management.ChangeOrgMemberRequest) *domain.Member {
	return domain.NewMember(authz.GetCtxData(ctx).OrgID, member.UserId, member.Roles...)
}

func orgMemberFromDomain(member *domain.Member) *management.OrgMember {
	return &management.OrgMember{
		UserId:     member.UserID,
		ChangeDate: timestamppb.New(member.ChangeDate),
		Roles:      member.Roles,
		Sequence:   member.Sequence,
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

func orgMemberSearchMethodToModel(key management.SearchMethod) domain.SearchMethod {
	switch key {
	case management.SearchMethod_SEARCHMETHOD_CONTAINS:
		return domain.SearchMethodContains
	case management.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return domain.SearchMethodContainsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_EQUALS:
		return domain.SearchMethodEquals
	case management.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return domain.SearchMethodEqualsIgnoreCase
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return domain.SearchMethodStartsWith
	case management.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return domain.SearchMethodStartsWithIgnoreCase
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

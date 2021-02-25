package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/domain"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func addIamMemberToDomain(member *admin.AddIamMemberRequest) *domain.Member {
	return domain.NewMember(domain.IAMID, member.UserId, member.Roles...)
}

func changeIamMemberToDomain(member *admin.ChangeIamMemberRequest) *domain.Member {
	return domain.NewMember(domain.IAMID, member.UserId, member.Roles...)
}

func iamMemberFromDomain(member *domain.Member) *admin.IamMember {
	return &admin.IamMember{
		UserId:     member.UserID,
		ChangeDate: timestamppb.New(member.ChangeDate),
		Roles:      member.Roles,
		Sequence:   member.Sequence,
	}
}

func iamMemberSearchRequestToModel(request *admin.IamMemberSearchRequest) *iam_model.IAMMemberSearchRequest {
	return &iam_model.IAMMemberSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: iamMemberSearchQueriesToModel(request.Queries),
	}
}

func iamMemberSearchQueriesToModel(queries []*admin.IamMemberSearchQuery) []*iam_model.IAMMemberSearchQuery {
	modelQueries := make([]*iam_model.IAMMemberSearchQuery, len(queries))
	for i, query := range queries {
		modelQueries[i] = iamMemberSearchQueryToModel(query)
	}

	return modelQueries
}

func iamMemberSearchQueryToModel(query *admin.IamMemberSearchQuery) *iam_model.IAMMemberSearchQuery {
	return &iam_model.IAMMemberSearchQuery{
		Key:    iamMemberSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func iamMemberSearchKeyToModel(key admin.IamMemberSearchKey) iam_model.IAMMemberSearchKey {
	switch key {
	case admin.IamMemberSearchKey_IAMMEMBERSEARCHKEY_EMAIL:
		return iam_model.IAMMemberSearchKeyEmail
	case admin.IamMemberSearchKey_IAMMEMBERSEARCHKEY_FIRST_NAME:
		return iam_model.IAMMemberSearchKeyFirstName
	case admin.IamMemberSearchKey_IAMMEMBERSEARCHKEY_LAST_NAME:
		return iam_model.IAMMemberSearchKeyLastName
	case admin.IamMemberSearchKey_IAMMEMBERSEARCHKEY_USER_ID:
		return iam_model.IAMMemberSearchKeyUserID
	default:
		return iam_model.IAMMemberSearchKeyUnspecified
	}
}

func searchMethodToModel(key admin.SearchMethod) domain.SearchMethod {
	switch key {
	case admin.SearchMethod_SEARCHMETHOD_CONTAINS:
		return domain.SearchMethodContains
	case admin.SearchMethod_SEARCHMETHOD_CONTAINS_IGNORE_CASE:
		return domain.SearchMethodContainsIgnoreCase
	case admin.SearchMethod_SEARCHMETHOD_EQUALS:
		return domain.SearchMethodEquals
	case admin.SearchMethod_SEARCHMETHOD_EQUALS_IGNORE_CASE:
		return domain.SearchMethodEqualsIgnoreCase
	case admin.SearchMethod_SEARCHMETHOD_STARTS_WITH:
		return domain.SearchMethodStartsWith
	case admin.SearchMethod_SEARCHMETHOD_STARTS_WITH_IGNORE_CASE:
		return domain.SearchMethodStartsWithIgnoreCase
	default:
		return -1
	}
}

func iamMemberSearchResponseFromModel(resp *iam_model.IAMMemberSearchResponse) *admin.IamMemberSearchResponse {
	timestamp, err := ptypes.TimestampProto(resp.Timestamp)
	logging.Log("GRPC-5shu8").OnError(err).Debug("date parse failed")
	return &admin.IamMemberSearchResponse{
		Limit:             resp.Limit,
		Offset:            resp.Offset,
		TotalResult:       resp.TotalResult,
		Result:            iamMembersFromView(resp.Result),
		ProcessedSequence: resp.Sequence,
		ViewTimestamp:     timestamp,
	}
}
func iamMembersFromView(viewMembers []*iam_model.IAMMemberView) []*admin.IamMemberView {
	members := make([]*admin.IamMemberView, len(viewMembers))

	for i, member := range viewMembers {
		members[i] = iamMemberFromView(member)
	}

	return members
}

func iamMemberFromView(member *iam_model.IAMMemberView) *admin.IamMemberView {
	changeDate, err := ptypes.TimestampProto(member.ChangeDate)
	logging.Log("GRPC-Lso9c").OnError(err).Debug("unable to parse changedate")
	creationDate, err := ptypes.TimestampProto(member.CreationDate)
	logging.Log("GRPC-6szE").OnError(err).Debug("unable to parse creation date")

	return &admin.IamMemberView{
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

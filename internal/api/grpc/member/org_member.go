package member

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	org_model "github.com/caos/zitadel/internal/org/model"
	member_pb "github.com/caos/zitadel/pkg/grpc/member"
)

func OrgMembersToPb(members []*org_model.OrgMemberView) []*member_pb.Member {
	m := make([]*member_pb.Member, len(members))
	for i, member := range members {
		m[i] = OrgMemberToPb(member)
	}
	return m
}

func OrgMemberToPb(m *org_model.OrgMemberView) *member_pb.Member {
	return &member_pb.Member{
		UserId:             m.UserID,
		Roles:              m.Roles,
		PreferredLoginName: m.PreferredLoginName,
		Email:              m.Email,
		FirstName:          m.FirstName,
		LastName:           m.LastName,
		DisplayName:        m.DisplayName,
		AvatarUrl:          m.AvatarURL,
		Details: object.ToViewDetailsPb(
			m.Sequence,
			m.CreationDate,
			m.ChangeDate,
			"", //TODO: not returnd
		),
	}
}

func MemberQueriesToOrgMember(queries []*member_pb.SearchQuery) []*org_model.OrgMemberSearchQuery {
	q := make([]*org_model.OrgMemberSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = MemberQueryToOrgMember(query)
	}
	return q
}

func MemberQueryToOrgMember(query *member_pb.SearchQuery) *org_model.OrgMemberSearchQuery {
	switch q := query.Query.(type) {
	case *member_pb.SearchQuery_EmailQuery:
		return EmailQueryToOrgMemberQuery(q.EmailQuery)
	case *member_pb.SearchQuery_FirstNameQuery:
		return FirstNameQueryToOrgMemberQuery(q.FirstNameQuery)
	case *member_pb.SearchQuery_LastNameQuery:
		return LastNameQueryToOrgMemberQuery(q.LastNameQuery)
	case *member_pb.SearchQuery_UserIdQuery:
		return UserIDQueryToOrgMemberQuery(q.UserIdQuery)
	default:
		return nil
	}
}

func FirstNameQueryToOrgMemberQuery(query *member_pb.FirstNameQuery) *org_model.OrgMemberSearchQuery {
	return &org_model.OrgMemberSearchQuery{
		Key:    org_model.OrgMemberSearchKeyFirstName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.FirstName,
	}
}

func LastNameQueryToOrgMemberQuery(query *member_pb.LastNameQuery) *org_model.OrgMemberSearchQuery {
	return &org_model.OrgMemberSearchQuery{
		Key:    org_model.OrgMemberSearchKeyLastName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.LastName,
	}
}

func EmailQueryToOrgMemberQuery(query *member_pb.EmailQuery) *org_model.OrgMemberSearchQuery {
	return &org_model.OrgMemberSearchQuery{
		Key:    org_model.OrgMemberSearchKeyEmail,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.Email,
	}
}

func UserIDQueryToOrgMemberQuery(query *member_pb.UserIDQuery) *org_model.OrgMemberSearchQuery {
	return &org_model.OrgMemberSearchQuery{
		Key:    org_model.OrgMemberSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  query.UserId,
	}
}

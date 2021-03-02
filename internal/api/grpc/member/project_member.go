package member

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	proj_model "github.com/caos/zitadel/internal/project/model"
	member_pb "github.com/caos/zitadel/pkg/grpc/member"
)

func ProjectMembersToPb(members []*proj_model.ProjectMemberView) []*member_pb.Member {
	m := make([]*member_pb.Member, len(members))
	for i, member := range members {
		m[i] = ProjectMemberToPb(member)
	}
	return m
}

func ProjectMemberToPb(m *proj_model.ProjectMemberView) *member_pb.Member {
	return &member_pb.Member{
		UserId: m.UserID,
		Roles:  m.Roles,
		// PreferredLoginName: //TODO: not implemented in be
		Email:       m.Email,
		FirstName:   m.FirstName,
		LastName:    m.LastName,
		DisplayName: m.DisplayName,
		Details: object.ToDetailsPb(
			m.Sequence,
			m.ChangeDate,
			"m.ResourceOwner", //TODO: not returnd
		),
	}
}

func MemberQueriesToProjectMember(queries []*member_pb.SearchQuery) []*proj_model.ProjectMemberSearchQuery {
	q := make([]*proj_model.ProjectMemberSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = MemberQueryToProjectMember(query)
	}
	return q
}

func MemberQueryToProjectMember(query *member_pb.SearchQuery) *proj_model.ProjectMemberSearchQuery {
	switch q := query.Query.(type) {
	case *member_pb.SearchQuery_Email:
		return EmailQueryToProjectMemberQuery(q.Email)
	case *member_pb.SearchQuery_FirstName:
		return FirstNameQueryToProjectMemberQuery(q.FirstName)
	case *member_pb.SearchQuery_LastName:
		return LastNameQueryToProjectMemberQuery(q.LastName)
	case *member_pb.SearchQuery_UserId:
		return UserIDQueryToProjectMemberQuery(q.UserId)
	default:
		return nil
	}
}

func FirstNameQueryToProjectMemberQuery(query *member_pb.FirstNameQuery) *proj_model.ProjectMemberSearchQuery {
	return &proj_model.ProjectMemberSearchQuery{
		Key:    proj_model.ProjectMemberSearchKeyFirstName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.FirstName,
	}
}

func LastNameQueryToProjectMemberQuery(query *member_pb.LastNameQuery) *proj_model.ProjectMemberSearchQuery {
	return &proj_model.ProjectMemberSearchQuery{
		Key:    proj_model.ProjectMemberSearchKeyLastName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.LastName,
	}
}

func EmailQueryToProjectMemberQuery(query *member_pb.EmailQuery) *proj_model.ProjectMemberSearchQuery {
	return &proj_model.ProjectMemberSearchQuery{
		Key:    proj_model.ProjectMemberSearchKeyEmail,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.Email,
	}
}

func UserIDQueryToProjectMemberQuery(query *member_pb.UserIDQuery) *proj_model.ProjectMemberSearchQuery {
	return &proj_model.ProjectMemberSearchQuery{
		Key:    proj_model.ProjectMemberSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  query.UserId,
	}
}

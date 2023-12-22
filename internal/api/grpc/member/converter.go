package member

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	member_pb "github.com/zitadel/zitadel/pkg/grpc/member"
)

func MemberToDomain(member *member_pb.Member) *domain.Member {
	return &domain.Member{
		UserID: member.UserId,
		Roles:  member.Roles,
	}
}

func MembersToPb(assetAPIPrefix string, members []*query.Member) []*member_pb.Member {
	m := make([]*member_pb.Member, len(members))
	for i, member := range members {
		m[i] = MemberToPb(assetAPIPrefix, member)
	}
	return m
}

func MemberToPb(assetAPIPrefix string, m *query.Member) *member_pb.Member {
	return &member_pb.Member{
		UserId:             m.UserID,
		Roles:              m.Roles,
		PreferredLoginName: m.PreferredLoginName,
		Email:              m.Email,
		FirstName:          m.FirstName,
		LastName:           m.LastName,
		DisplayName:        m.DisplayName,
		AvatarUrl:          domain.AvatarURL(assetAPIPrefix, m.ResourceOwner, m.AvatarURL),
		UserType:           user.TypeToPb(m.UserType),
		Details: object.ToViewDetailsPb(
			m.Sequence,
			m.CreationDate,
			m.ChangeDate,
			m.ResourceOwner,
		),
	}
}

func MemberQueriesToQuery(queries []*member_pb.SearchQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = MemberQueryToMember(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func MemberQueryToMember(search *member_pb.SearchQuery) (query.SearchQuery, error) {
	switch q := search.Query.(type) {
	case *member_pb.SearchQuery_EmailQuery:
		return query.NewMemberEmailSearchQuery(object.TextMethodToQuery(q.EmailQuery.Method), q.EmailQuery.Email)
	case *member_pb.SearchQuery_FirstNameQuery:
		return query.NewMemberFirstNameSearchQuery(object.TextMethodToQuery(q.FirstNameQuery.Method), q.FirstNameQuery.FirstName)
	case *member_pb.SearchQuery_LastNameQuery:
		return query.NewMemberLastNameSearchQuery(object.TextMethodToQuery(q.LastNameQuery.Method), q.LastNameQuery.LastName)
	case *member_pb.SearchQuery_UserIdQuery:
		return query.NewMemberUserIDSearchQuery(q.UserIdQuery.UserId)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "MEMBE-7Bb92", "Errors.Query.InvalidRequest")
	}
}

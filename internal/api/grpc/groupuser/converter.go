package groupuser

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	groupuser_pb "github.com/zitadel/zitadel/pkg/grpc/groupusers"
)

func MembersToPb(assetAPIPrefix string, users []*query.GroupUser) []*groupuser_pb.GroupUser {
	m := make([]*groupuser_pb.GroupUser, len(users))
	for i, user := range users {
		m[i] = MemberToPb(assetAPIPrefix, user)
	}
	return m
}

func MemberToPb(assetAPIPrefix string, m *query.GroupUser) *groupuser_pb.GroupUser {
	return &groupuser_pb.GroupUser{
		UserId:             m.UserID,
		GroupId:            m.GroupID,
		Attributes:         m.Attributes,
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

func MemberQueriesToQuery(queries []*groupuser_pb.SearchQuery) (q []query.SearchQuery, err error) {
	q = make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = MemberQueryToMember(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func MemberQueryToMember(search *groupuser_pb.SearchQuery) (query.SearchQuery, error) {
	switch q := search.Query.(type) {
	case *groupuser_pb.SearchQuery_EmailQuery:
		return query.NewGroupUserEmailSearchQuery(object.TextMethodToQuery(q.EmailQuery.Method), q.EmailQuery.Email)
	case *groupuser_pb.SearchQuery_FirstNameQuery:
		return query.NewGroupUserFirstNameSearchQuery(object.TextMethodToQuery(q.FirstNameQuery.Method), q.FirstNameQuery.FirstName)
	case *groupuser_pb.SearchQuery_LastNameQuery:
		return query.NewGroupUserLastNameSearchQuery(object.TextMethodToQuery(q.LastNameQuery.Method), q.LastNameQuery.LastName)
	case *groupuser_pb.SearchQuery_UserIdQuery:
		return query.NewGroupUserUserIDSearchQuery(q.UserIdQuery.UserId)
	case *groupuser_pb.SearchQuery_AttributesQuery:
		return query.NewGroupUserAttributesSearchQuery(q.AttributesQuery.Attributes)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GMEMBE-8Bb92", "Errors.Query.InvalidRequest")
	}
}

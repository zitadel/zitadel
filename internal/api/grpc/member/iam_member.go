package member

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	member_pb "github.com/caos/zitadel/pkg/grpc/member"
)

func IAMMembersToPb(members []*iam_model.IAMMemberView) []*member_pb.Member {
	m := make([]*member_pb.Member, len(members))
	for i, member := range members {
		m[i] = IAMMemberToPb(member)
	}
	return m
}

func IAMMemberToPb(m *iam_model.IAMMemberView) *member_pb.Member {
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

func MemberQueriesToIAMMember(queries []*member_pb.SearchQuery) []*iam_model.IAMMemberSearchQuery {
	q := make([]*iam_model.IAMMemberSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = MemberQueryToIAMMember(query)
	}
	return q
}

func MemberQueryToIAMMember(query *member_pb.SearchQuery) *iam_model.IAMMemberSearchQuery {
	switch q := query.Query.(type) {
	case *member_pb.SearchQuery_EmailQuery:
		return EmailQueryToIAMMemberQuery(q.EmailQuery)
	case *member_pb.SearchQuery_FirstNameQuery:
		return FirstNameQueryToIAMMemberQuery(q.FirstNameQuery)
	case *member_pb.SearchQuery_LastNameQuery:
		return LastNameQueryToIAMMemberQuery(q.LastNameQuery)
	case *member_pb.SearchQuery_UserIdQuery:
		return UserIDQueryToIAMMemberQuery(q.UserIdQuery)
	default:
		return nil
	}
}

func FirstNameQueryToIAMMemberQuery(query *member_pb.FirstNameQuery) *iam_model.IAMMemberSearchQuery {
	return &iam_model.IAMMemberSearchQuery{
		Key:    iam_model.IAMMemberSearchKeyFirstName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.FirstName,
	}
}

func LastNameQueryToIAMMemberQuery(query *member_pb.LastNameQuery) *iam_model.IAMMemberSearchQuery {
	return &iam_model.IAMMemberSearchQuery{
		Key:    iam_model.IAMMemberSearchKeyLastName,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.LastName,
	}
}

func EmailQueryToIAMMemberQuery(query *member_pb.EmailQuery) *iam_model.IAMMemberSearchQuery {
	return &iam_model.IAMMemberSearchQuery{
		Key:    iam_model.IAMMemberSearchKeyEmail,
		Method: object.TextMethodToModel(query.Method),
		Value:  query.Email,
	}
}

func UserIDQueryToIAMMemberQuery(query *member_pb.UserIDQuery) *iam_model.IAMMemberSearchQuery {
	return &iam_model.IAMMemberSearchQuery{
		Key:    iam_model.IAMMemberSearchKeyUserID,
		Method: domain.SearchMethodEquals,
		Value:  query.UserId,
	}
}

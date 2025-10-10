package user

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
)

func UserQueriesToQuery(queries []*user_pb.SearchQuery, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = UserQueryToQuery(query, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func UserQueryToQuery(query *user_pb.SearchQuery, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := query.Query.(type) {
	case *user_pb.SearchQuery_UserNameQuery:
		return UserNameQueryToQuery(q.UserNameQuery)
	case *user_pb.SearchQuery_FirstNameQuery:
		return FirstNameQueryToQuery(q.FirstNameQuery)
	case *user_pb.SearchQuery_LastNameQuery:
		return LastNameQueryToQuery(q.LastNameQuery)
	case *user_pb.SearchQuery_NickNameQuery:
		return NickNameQueryToQuery(q.NickNameQuery)
	case *user_pb.SearchQuery_DisplayNameQuery:
		return DisplayNameQueryToQuery(q.DisplayNameQuery)
	case *user_pb.SearchQuery_EmailQuery:
		return EmailQueryToQuery(q.EmailQuery)
	case *user_pb.SearchQuery_StateQuery:
		return StateQueryToQuery(q.StateQuery)
	case *user_pb.SearchQuery_TypeQuery:
		return TypeQueryToQuery(q.TypeQuery)
	case *user_pb.SearchQuery_LoginNameQuery:
		return LoginNameQueryToQuery(q.LoginNameQuery)
	case *user_pb.SearchQuery_ResourceOwner:
		return ResourceOwnerQueryToQuery(q.ResourceOwner)
	case *user_pb.SearchQuery_InUserIdsQuery:
		return InUserIdsQueryToQuery(q.InUserIdsQuery)
	case *user_pb.SearchQuery_InUserEmailsQuery:
		return InUserEmailsQueryToQuery(q.InUserEmailsQuery)
	case *user_pb.SearchQuery_OrQuery:
		return OrQueryToQuery(q.OrQuery, level)
	case *user_pb.SearchQuery_AndQuery:
		return AndQueryToQuery(q.AndQuery, level)
	case *user_pb.SearchQuery_NotQuery:
		return NotQueryToQuery(q.NotQuery, level)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func UserNameQueryToQuery(q *user_pb.UserNameQuery) (query.SearchQuery, error) {
	return query.NewUserUsernameSearchQuery(q.UserName, object.TextMethodToQuery(q.Method))
}

func FirstNameQueryToQuery(q *user_pb.FirstNameQuery) (query.SearchQuery, error) {
	return query.NewUserFirstNameSearchQuery(q.FirstName, object.TextMethodToQuery(q.Method))
}

func LastNameQueryToQuery(q *user_pb.LastNameQuery) (query.SearchQuery, error) {
	return query.NewUserLastNameSearchQuery(q.LastName, object.TextMethodToQuery(q.Method))
}

func NickNameQueryToQuery(q *user_pb.NickNameQuery) (query.SearchQuery, error) {
	return query.NewUserNickNameSearchQuery(q.NickName, object.TextMethodToQuery(q.Method))
}

func DisplayNameQueryToQuery(q *user_pb.DisplayNameQuery) (query.SearchQuery, error) {
	return query.NewUserDisplayNameSearchQuery(q.DisplayName, object.TextMethodToQuery(q.Method))
}

func EmailQueryToQuery(q *user_pb.EmailQuery) (query.SearchQuery, error) {
	return query.NewUserEmailSearchQuery(q.EmailAddress, object.TextMethodToQuery(q.Method))
}

func StateQueryToQuery(q *user_pb.StateQuery) (query.SearchQuery, error) {
	return query.NewUserStateSearchQuery(q.State.ToDomain())
}

func TypeQueryToQuery(q *user_pb.TypeQuery) (query.SearchQuery, error) {
	return query.NewUserTypeSearchQuery(q.Type.ToDomain())
}

func LoginNameQueryToQuery(q *user_pb.LoginNameQuery) (query.SearchQuery, error) {
	return query.NewUserLoginNameExistsQuery(q.LoginName, object.TextMethodToQuery(q.Method))
}

func ResourceOwnerQueryToQuery(q *user_pb.ResourceOwnerQuery) (query.SearchQuery, error) {
	return query.NewUserResourceOwnerSearchQuery(q.OrgID, query.TextEquals)
}

func InUserIdsQueryToQuery(q *user_pb.InUserIDQuery) (query.SearchQuery, error) {
	return query.NewUserInUserIdsSearchQuery(q.UserIds)
}

func InUserEmailsQueryToQuery(q *user_pb.InUserEmailsQuery) (query.SearchQuery, error) {
	return query.NewUserInUserEmailsSearchQuery(q.UserEmails)
}

func OrQueryToQuery(q *user_pb.OrQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := UserQueriesToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserOrSearchQuery(mappedQueries)
}
func AndQueryToQuery(q *user_pb.AndQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := UserQueriesToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserAndSearchQuery(mappedQueries)
}
func NotQueryToQuery(q *user_pb.NotQuery, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := UserQueryToQuery(q.Query, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserNotSearchQuery(mappedQuery)
}

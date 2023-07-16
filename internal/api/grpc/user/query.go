package user

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
)

func UserQueriesToQuery(queries []*user_pb.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = UserQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func UserQueryToQuery(query *user_pb.SearchQuery) (query.SearchQuery, error) {
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
	default:
		return nil, errors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
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
	return query.NewUserStateSearchQuery(int32(q.State))
}

func TypeQueryToQuery(q *user_pb.TypeQuery) (query.SearchQuery, error) {
	return query.NewUserTypeSearchQuery(int32(q.Type))
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

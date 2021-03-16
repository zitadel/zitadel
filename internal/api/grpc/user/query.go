package user

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	user_model "github.com/caos/zitadel/internal/user/model"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func UserQueriesToModel(queries []*user_pb.SearchQuery) []*user_model.UserSearchQuery {
	q := make([]*user_model.UserSearchQuery, len(queries))
	for i, query := range queries {
		q[i] = UserQueryToModel(query)
	}
	return q
}

func UserQueryToModel(query *user_pb.SearchQuery) *user_model.UserSearchQuery {
	switch q := query.Query.(type) {
	case *user_pb.SearchQuery_UserNameQuery:
		return UserNameQueryToModel(q.UserNameQuery)
	case *user_pb.SearchQuery_FirstNameQuery:
		return FirstNameQueryToModel(q.FirstNameQuery)
	case *user_pb.SearchQuery_LastNameQuery:
		return LastNameQueryToModel(q.LastNameQuery)
	case *user_pb.SearchQuery_NickNameQuery:
		return NickNameQueryToModel(q.NickNameQuery)
	case *user_pb.SearchQuery_DisplayNameQuery:
		return DisplayNameQueryToModel(q.DisplayNameQuery)
	case *user_pb.SearchQuery_EmailQuery:
		return EmailQueryToModel(q.EmailQuery)
	case *user_pb.SearchQuery_StateQuery:
		return StateQueryToModel(q.StateQuery)
	case *user_pb.SearchQuery_TypeQuery:
		return TypeQueryToModel(q.TypeQuery)
	case *user_pb.SearchQuery_ResourceOwner:
		return ResourceOwnerQueryToModel(q.ResourceOwner)
	default:
		return nil
	}
}

func UserNameQueryToModel(q *user_pb.UserNameQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyUserName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.UserName,
	}
}

func FirstNameQueryToModel(q *user_pb.FirstNameQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyFirstName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.FirstName,
	}
}

func LastNameQueryToModel(q *user_pb.LastNameQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyLastName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.LastName,
	}
}

func NickNameQueryToModel(q *user_pb.NickNameQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyNickName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.NickName,
	}
}

func DisplayNameQueryToModel(q *user_pb.DisplayNameQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyDisplayName,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.DisplayName,
	}
}

func EmailQueryToModel(q *user_pb.EmailQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyEmail,
		Method: object.TextMethodToModel(q.Method),
		Value:  q.EmailAddress,
	}
}

func StateQueryToModel(q *user_pb.StateQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyState,
		Method: domain.SearchMethodEquals,
		Value:  q.State,
	}
}

func TypeQueryToModel(q *user_pb.TypeQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyType,
		Method: domain.SearchMethodEquals,
		Value:  UserTypeToModel(q.Type),
	}
}

func UserTypeToModel(t user_pb.Type) string {
	switch t {
	case user_pb.Type_TYPE_HUMAN:
		return "human"
	case user_pb.Type_TYPE_MACHINE:
		return "machine"
	default:
		return ""
	}
}

func ResourceOwnerQueryToModel(q *user_pb.ResourceOwnerQuery) *user_model.UserSearchQuery {
	return &user_model.UserSearchQuery{
		Key:    user_model.UserSearchKeyResourceOwner,
		Method: domain.SearchMethodEquals,
		Value:  q.OrgID,
	}
}

package user

import (
	"context"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func (s *Server) GetUserByID(ctx context.Context, req *user.GetUserByIDRequest) (_ *user.GetUserByIDResponse, err error) {
	resp, err := s.query.GetUserByID(ctx, true, req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &user.GetUserByIDResponse{
		Details: object.DomainToDetailsPb(&domain.ObjectDetails{
			Sequence:      resp.Sequence,
			EventDate:     resp.ChangeDate,
			ResourceOwner: resp.ResourceOwner,
		}),
		User: UserToPb(resp, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	queries, err := ListUsersRequestToModel(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUsers(ctx, queries)
	if err != nil {
		return nil, err
	}
	res.RemoveNoPermission(ctx, s.query)
	return &user.ListUsersResponse{
		Result:  UsersToPb(res.Users, s.assetAPIPrefix(ctx)),
		Details: object.ToListDetails(res.SearchResponse),
	}, nil
}

func UsersToPb(users []*query.User, assetPrefix string) []*user.User {
	u := make([]*user.User, len(users))
	for i, user := range users {
		u[i] = UserToPb(user, assetPrefix)
	}
	return u
}

func UserToPb(userQ *query.User, assetPrefix string) *user.User {
	return &user.User{
		UserId:             userQ.ID,
		State:              UserStateToPb(userQ.State),
		Username:           userQ.Username,
		LoginNames:         userQ.LoginNames,
		PreferredLoginName: userQ.PreferredLoginName,
		Type:               UserTypeToPb(userQ, assetPrefix),
	}
}

func UserTypeToPb(userQ *query.User, assetPrefix string) user.UserType {
	if userQ.Human != nil {
		return &user.User_Human{
			Human: HumanToPb(userQ.Human, assetPrefix, userQ.ResourceOwner),
		}
	}
	if userQ.Machine != nil {
		return &user.User_Machine{
			Machine: MachineToPb(userQ.Machine),
		}
	}
	return nil
}

func HumanToPb(userQ *query.Human, assetPrefix, owner string) *user.HumanUser {
	return &user.HumanUser{
		Profile: &user.HumanProfile{
			GivenName:         userQ.FirstName,
			FamilyName:        userQ.LastName,
			NickName:          gu.Ptr(userQ.NickName),
			DisplayName:       gu.Ptr(userQ.DisplayName),
			PreferredLanguage: gu.Ptr(userQ.PreferredLanguage.String()),
			Gender:            gu.Ptr(GenderToPb(userQ.Gender)),
			AvatarUrl:         domain.AvatarURL(assetPrefix, owner, userQ.AvatarKey),
		},
		Email: &user.HumanEmail{
			Email:      string(userQ.Email),
			IsVerified: userQ.IsEmailVerified,
		},
		Phone: &user.HumanPhone{
			Phone:      string(userQ.Phone),
			IsVerified: userQ.IsPhoneVerified,
		},
	}
}

func MachineToPb(userQ *query.Machine) *user.MachineUser {
	return &user.MachineUser{
		Name:            userQ.Name,
		Description:     userQ.Description,
		HasSecret:       userQ.Secret != nil,
		AccessTokenType: AccessTokenTypeToPb(userQ.AccessTokenType),
	}
}

func UserStateToPb(state domain.UserState) user.UserState {
	switch state {
	case domain.UserStateActive:
		return user.UserState_USER_STATE_ACTIVE
	case domain.UserStateInactive:
		return user.UserState_USER_STATE_INACTIVE
	case domain.UserStateDeleted:
		return user.UserState_USER_STATE_DELETED
	case domain.UserStateInitial:
		return user.UserState_USER_STATE_INITIAL
	case domain.UserStateLocked:
		return user.UserState_USER_STATE_LOCKED
	default:
		return user.UserState_USER_STATE_UNSPECIFIED
	}
}

func GenderToPb(gender domain.Gender) user.Gender {
	switch gender {
	case domain.GenderDiverse:
		return user.Gender_GENDER_DIVERSE
	case domain.GenderFemale:
		return user.Gender_GENDER_FEMALE
	case domain.GenderMale:
		return user.Gender_GENDER_MALE
	default:
		return user.Gender_GENDER_UNSPECIFIED
	}
}

func AccessTokenTypeToPb(accessTokenType domain.OIDCTokenType) user.AccessTokenType {
	switch accessTokenType {
	case domain.OIDCTokenTypeBearer:
		return user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER
	case domain.OIDCTokenTypeJWT:
		return user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT
	default:
		return user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER
	}
}

func ListUsersRequestToModel(req *user.ListUsersRequest) (*query.UserSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := UserQueriesToQuery(req.Queries, 0 /*start from level 0*/)
	if err != nil {
		return nil, err
	}
	return &query.UserSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: UserFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func UserFieldNameToSortingColumn(field user.UserFieldName) query.Column {
	switch field {
	case user.UserFieldName_USER_FIELD_NAME_EMAIL:
		return query.HumanEmailCol
	case user.UserFieldName_USER_FIELD_NAME_FIRST_NAME:
		return query.HumanFirstNameCol
	case user.UserFieldName_USER_FIELD_NAME_LAST_NAME:
		return query.HumanLastNameCol
	case user.UserFieldName_USER_FIELD_NAME_DISPLAY_NAME:
		return query.HumanDisplayNameCol
	case user.UserFieldName_USER_FIELD_NAME_USER_NAME:
		return query.UserUsernameCol
	case user.UserFieldName_USER_FIELD_NAME_STATE:
		return query.UserStateCol
	case user.UserFieldName_USER_FIELD_NAME_TYPE:
		return query.UserTypeCol
	case user.UserFieldName_USER_FIELD_NAME_NICK_NAME:
		return query.HumanNickNameCol
	case user.UserFieldName_USER_FIELD_NAME_CREATION_DATE:
		return query.UserCreationDateCol
	default:
		return query.UserIDCol
	}
}

func UserQueriesToQuery(queries []*user.SearchQuery, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = UserQueryToQuery(query, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func UserQueryToQuery(query *user.SearchQuery, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-zsQ97", "Errors.User.TooManyNestingLevels")
	}
	switch q := query.Query.(type) {
	case *user.SearchQuery_UserNameQuery:
		return UserNameQueryToQuery(q.UserNameQuery)
	case *user.SearchQuery_FirstNameQuery:
		return FirstNameQueryToQuery(q.FirstNameQuery)
	case *user.SearchQuery_LastNameQuery:
		return LastNameQueryToQuery(q.LastNameQuery)
	case *user.SearchQuery_NickNameQuery:
		return NickNameQueryToQuery(q.NickNameQuery)
	case *user.SearchQuery_DisplayNameQuery:
		return DisplayNameQueryToQuery(q.DisplayNameQuery)
	case *user.SearchQuery_EmailQuery:
		return EmailQueryToQuery(q.EmailQuery)
	case *user.SearchQuery_StateQuery:
		return StateQueryToQuery(q.StateQuery)
	case *user.SearchQuery_TypeQuery:
		return TypeQueryToQuery(q.TypeQuery)
	case *user.SearchQuery_LoginNameQuery:
		return LoginNameQueryToQuery(q.LoginNameQuery)
	case *user.SearchQuery_ResourceOwner:
		return ResourceOwnerQueryToQuery(q.ResourceOwner)
	case *user.SearchQuery_InUserIdsQuery:
		return InUserIdsQueryToQuery(q.InUserIdsQuery)
	case *user.SearchQuery_OrQuery:
		return OrQueryToQuery(q.OrQuery, level)
	case *user.SearchQuery_AndQuery:
		return AndQueryToQuery(q.AndQuery, level)
	case *user.SearchQuery_NotQuery:
		return NotQueryToQuery(q.NotQuery, level)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func UserNameQueryToQuery(q *user.UserNameQuery) (query.SearchQuery, error) {
	return query.NewUserUsernameSearchQuery(q.UserName, object.TextMethodToQuery(q.Method))
}

func FirstNameQueryToQuery(q *user.FirstNameQuery) (query.SearchQuery, error) {
	return query.NewUserFirstNameSearchQuery(q.FirstName, object.TextMethodToQuery(q.Method))
}

func LastNameQueryToQuery(q *user.LastNameQuery) (query.SearchQuery, error) {
	return query.NewUserLastNameSearchQuery(q.LastName, object.TextMethodToQuery(q.Method))
}

func NickNameQueryToQuery(q *user.NickNameQuery) (query.SearchQuery, error) {
	return query.NewUserNickNameSearchQuery(q.NickName, object.TextMethodToQuery(q.Method))
}

func DisplayNameQueryToQuery(q *user.DisplayNameQuery) (query.SearchQuery, error) {
	return query.NewUserDisplayNameSearchQuery(q.DisplayName, object.TextMethodToQuery(q.Method))
}

func EmailQueryToQuery(q *user.EmailQuery) (query.SearchQuery, error) {
	return query.NewUserEmailSearchQuery(q.EmailAddress, object.TextMethodToQuery(q.Method))
}

func StateQueryToQuery(q *user.StateQuery) (query.SearchQuery, error) {
	return query.NewUserStateSearchQuery(int32(q.State))
}

func TypeQueryToQuery(q *user.TypeQuery) (query.SearchQuery, error) {
	return query.NewUserTypeSearchQuery(int32(q.Type))
}

func LoginNameQueryToQuery(q *user.LoginNameQuery) (query.SearchQuery, error) {
	return query.NewUserLoginNameExistsQuery(q.LoginName, object.TextMethodToQuery(q.Method))
}

func ResourceOwnerQueryToQuery(q *user.ResourceOwnerQuery) (query.SearchQuery, error) {
	return query.NewUserResourceOwnerSearchQuery(q.OrgID, query.TextEquals)
}

func InUserIdsQueryToQuery(q *user.InUserIDQuery) (query.SearchQuery, error) {
	return query.NewUserInUserIdsSearchQuery(q.UserIds)
}
func OrQueryToQuery(q *user.OrQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := UserQueriesToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserOrSearchQuery(mappedQueries)
}
func AndQueryToQuery(q *user.AndQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := UserQueriesToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserAndSearchQuery(mappedQueries)
}
func NotQueryToQuery(q *user.NotQuery, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := UserQueryToQuery(q.Query, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserNotSearchQuery(mappedQuery)
}

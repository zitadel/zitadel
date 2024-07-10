package user

import (
	"context"

	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (s *Server) GetUserByID(ctx context.Context, req *user.GetUserByIDRequest) (_ *user.GetUserByIDResponse, err error) {
	resp, err := s.query.GetUserByID(ctx, true, req.GetUserId())
	if err != nil {
		return nil, err
	}
	if authz.GetCtxData(ctx).UserID != req.GetUserId() {
		if err := s.checkPermission(ctx, domain.PermissionUserRead, resp.ResourceOwner, req.GetUserId()); err != nil {
			return nil, err
		}
	}
	return &user.GetUserByIDResponse{
		Details: object.DomainToDetailsPb(&domain.ObjectDetails{
			Sequence:      resp.Sequence,
			EventDate:     resp.ChangeDate,
			ResourceOwner: resp.ResourceOwner,
		}),
		User: userToPb(resp, s.assetAPIPrefix(ctx)),
	}, nil
}

func (s *Server) ListUsers(ctx context.Context, req *user.ListUsersRequest) (*user.ListUsersResponse, error) {
	queries, err := listUsersRequestToModel(req)
	if err != nil {
		return nil, err
	}
	res, err := s.query.SearchUsers(ctx, queries)
	if err != nil {
		return nil, err
	}
	res.RemoveNoPermission(ctx, s.checkPermission)
	return &user.ListUsersResponse{
		Result:  UsersToPb(res.Users, s.assetAPIPrefix(ctx)),
		Details: object.ToListDetails(res.SearchResponse),
	}, nil
}

func UsersToPb(users []*query.User, assetPrefix string) []*user.User {
	u := make([]*user.User, len(users))
	for i, user := range users {
		u[i] = userToPb(user, assetPrefix)
	}
	return u
}

func userToPb(userQ *query.User, assetPrefix string) *user.User {
	return &user.User{
		UserId: userQ.ID,
		Details: object.DomainToDetailsPb(&domain.ObjectDetails{
			Sequence:      userQ.Sequence,
			EventDate:     userQ.ChangeDate,
			ResourceOwner: userQ.ResourceOwner,
		}),
		State:              userStateToPb(userQ.State),
		Username:           userQ.Username,
		LoginNames:         userQ.LoginNames,
		PreferredLoginName: userQ.PreferredLoginName,
		Type:               userTypeToPb(userQ, assetPrefix),
	}
}

func userTypeToPb(userQ *query.User, assetPrefix string) user.UserType {
	if userQ.Human != nil {
		return &user.User_Human{
			Human: humanToPb(userQ.Human, assetPrefix, userQ.ResourceOwner),
		}
	}
	if userQ.Machine != nil {
		return &user.User_Machine{
			Machine: machineToPb(userQ.Machine),
		}
	}
	return nil
}

func humanToPb(userQ *query.Human, assetPrefix, owner string) *user.HumanUser {
	var passwordChanged *timestamppb.Timestamp
	if !userQ.PasswordChanged.IsZero() {
		passwordChanged = timestamppb.New(userQ.PasswordChanged)
	}
	return &user.HumanUser{
		Profile: &user.HumanProfile{
			GivenName:         userQ.FirstName,
			FamilyName:        userQ.LastName,
			NickName:          gu.Ptr(userQ.NickName),
			DisplayName:       gu.Ptr(userQ.DisplayName),
			PreferredLanguage: gu.Ptr(userQ.PreferredLanguage.String()),
			Gender:            gu.Ptr(genderToPb(userQ.Gender)),
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
		PasswordChangeRequired: userQ.PasswordChangeRequired,
		PasswordChanged:        passwordChanged,
	}
}

func machineToPb(userQ *query.Machine) *user.MachineUser {
	return &user.MachineUser{
		Name:            userQ.Name,
		Description:     userQ.Description,
		HasSecret:       userQ.EncodedSecret != "",
		AccessTokenType: accessTokenTypeToPb(userQ.AccessTokenType),
	}
}

func userStateToPb(state domain.UserState) user.UserState {
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
	case domain.UserStateUnspecified:
		return user.UserState_USER_STATE_UNSPECIFIED
	case domain.UserStateSuspend:
		return user.UserState_USER_STATE_UNSPECIFIED
	default:
		return user.UserState_USER_STATE_UNSPECIFIED
	}
}

func genderToPb(gender domain.Gender) user.Gender {
	switch gender {
	case domain.GenderDiverse:
		return user.Gender_GENDER_DIVERSE
	case domain.GenderFemale:
		return user.Gender_GENDER_FEMALE
	case domain.GenderMale:
		return user.Gender_GENDER_MALE
	case domain.GenderUnspecified:
		return user.Gender_GENDER_UNSPECIFIED
	default:
		return user.Gender_GENDER_UNSPECIFIED
	}
}

func accessTokenTypeToPb(accessTokenType domain.OIDCTokenType) user.AccessTokenType {
	switch accessTokenType {
	case domain.OIDCTokenTypeBearer:
		return user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER
	case domain.OIDCTokenTypeJWT:
		return user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT
	default:
		return user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER
	}
}

func listUsersRequestToModel(req *user.ListUsersRequest) (*query.UserSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := userQueriesToQuery(req.Queries, 0 /*start from level 0*/)
	if err != nil {
		return nil, err
	}
	return &query.UserSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: userFieldNameToSortingColumn(req.SortingColumn),
		},
		Queries: queries,
	}, nil
}

func userFieldNameToSortingColumn(field user.UserFieldName) query.Column {
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
	case user.UserFieldName_USER_FIELD_NAME_UNSPECIFIED:
		return query.UserIDCol
	default:
		return query.UserIDCol
	}
}

func userQueriesToQuery(queries []*user.SearchQuery, level uint8) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = userQueryToQuery(query, level)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func userQueryToQuery(query *user.SearchQuery, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := query.Query.(type) {
	case *user.SearchQuery_UserNameQuery:
		return userNameQueryToQuery(q.UserNameQuery)
	case *user.SearchQuery_FirstNameQuery:
		return firstNameQueryToQuery(q.FirstNameQuery)
	case *user.SearchQuery_LastNameQuery:
		return lastNameQueryToQuery(q.LastNameQuery)
	case *user.SearchQuery_NickNameQuery:
		return nickNameQueryToQuery(q.NickNameQuery)
	case *user.SearchQuery_DisplayNameQuery:
		return displayNameQueryToQuery(q.DisplayNameQuery)
	case *user.SearchQuery_EmailQuery:
		return emailQueryToQuery(q.EmailQuery)
	case *user.SearchQuery_StateQuery:
		return stateQueryToQuery(q.StateQuery)
	case *user.SearchQuery_TypeQuery:
		return typeQueryToQuery(q.TypeQuery)
	case *user.SearchQuery_LoginNameQuery:
		return loginNameQueryToQuery(q.LoginNameQuery)
	case *user.SearchQuery_OrganizationIdQuery:
		return resourceOwnerQueryToQuery(q.OrganizationIdQuery)
	case *user.SearchQuery_InUserIdsQuery:
		return inUserIdsQueryToQuery(q.InUserIdsQuery)
	case *user.SearchQuery_OrQuery:
		return orQueryToQuery(q.OrQuery, level)
	case *user.SearchQuery_AndQuery:
		return andQueryToQuery(q.AndQuery, level)
	case *user.SearchQuery_NotQuery:
		return notQueryToQuery(q.NotQuery, level)
	case *user.SearchQuery_InUserEmailsQuery:
		return inUserEmailsQueryToQuery(q.InUserEmailsQuery)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

func userNameQueryToQuery(q *user.UserNameQuery) (query.SearchQuery, error) {
	return query.NewUserUsernameSearchQuery(q.UserName, object.TextMethodToQuery(q.Method))
}

func firstNameQueryToQuery(q *user.FirstNameQuery) (query.SearchQuery, error) {
	return query.NewUserFirstNameSearchQuery(q.FirstName, object.TextMethodToQuery(q.Method))
}

func lastNameQueryToQuery(q *user.LastNameQuery) (query.SearchQuery, error) {
	return query.NewUserLastNameSearchQuery(q.LastName, object.TextMethodToQuery(q.Method))
}

func nickNameQueryToQuery(q *user.NickNameQuery) (query.SearchQuery, error) {
	return query.NewUserNickNameSearchQuery(q.NickName, object.TextMethodToQuery(q.Method))
}

func displayNameQueryToQuery(q *user.DisplayNameQuery) (query.SearchQuery, error) {
	return query.NewUserDisplayNameSearchQuery(q.DisplayName, object.TextMethodToQuery(q.Method))
}

func emailQueryToQuery(q *user.EmailQuery) (query.SearchQuery, error) {
	return query.NewUserEmailSearchQuery(q.EmailAddress, object.TextMethodToQuery(q.Method))
}

func stateQueryToQuery(q *user.StateQuery) (query.SearchQuery, error) {
	return query.NewUserStateSearchQuery(int32(q.State))
}

func typeQueryToQuery(q *user.TypeQuery) (query.SearchQuery, error) {
	return query.NewUserTypeSearchQuery(int32(q.Type))
}

func loginNameQueryToQuery(q *user.LoginNameQuery) (query.SearchQuery, error) {
	return query.NewUserLoginNameExistsQuery(q.LoginName, object.TextMethodToQuery(q.Method))
}

func resourceOwnerQueryToQuery(q *user.OrganizationIdQuery) (query.SearchQuery, error) {
	return query.NewUserResourceOwnerSearchQuery(q.OrganizationId, query.TextEquals)
}

func inUserIdsQueryToQuery(q *user.InUserIDQuery) (query.SearchQuery, error) {
	return query.NewUserInUserIdsSearchQuery(q.UserIds)
}
func orQueryToQuery(q *user.OrQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := userQueriesToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserOrSearchQuery(mappedQueries)
}
func andQueryToQuery(q *user.AndQuery, level uint8) (query.SearchQuery, error) {
	mappedQueries, err := userQueriesToQuery(q.Queries, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserAndSearchQuery(mappedQueries)
}
func notQueryToQuery(q *user.NotQuery, level uint8) (query.SearchQuery, error) {
	mappedQuery, err := userQueryToQuery(q.Query, level+1)
	if err != nil {
		return nil, err
	}
	return query.NewUserNotSearchQuery(mappedQuery)
}

func inUserEmailsQueryToQuery(q *user.InUserEmailsQuery) (query.SearchQuery, error) {
	return query.NewUserInUserEmailsSearchQuery(q.UserEmails)
}

package convert

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func ListUsersRequestToModel(req *user.ListUsersRequest) (*query.UserSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.GetQuery())
	queries, err := userQueriesToQuery(req.GetQueries(), 0 /*start from level 0*/)
	if err != nil {
		return nil, err
	}
	return &query.UserSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: userFieldNameToSortingColumn(req.GetSortingColumn()),
		},
		Queries: queries,
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
	if userQ == nil {
		return nil
	}

	return &user.User{
		UserId: userQ.ID,
		Details: object.DomainToDetailsPb(&domain.ObjectDetails{
			Sequence:      userQ.Sequence,
			EventDate:     userQ.ChangeDate,
			ResourceOwner: userQ.ResourceOwner,
			CreationDate:  userQ.CreationDate,
		}),
		State:              userStateToPb(userQ.State),
		Username:           userQ.Username,
		LoginNames:         userQ.LoginNames,
		PreferredLoginName: userQ.PreferredLoginName,
		Type:               userTypeToPb(userQ, assetPrefix),
	}
}

func userTypeToPb(userQ *query.User, assetPrefix string) user.UserType {
	if userQ == nil {
		return nil
	}

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

func userQueryToQuery(sq *user.SearchQuery, level uint8) (query.SearchQuery, error) {
	if level > 20 {
		// can't go deeper than 20 levels of nesting.
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-zsQ97", "Errors.Query.TooManyNestingLevels")
	}
	switch q := sq.Query.(type) {
	case *user.SearchQuery_UserNameQuery:
		return query.NewUserUsernameSearchQuery(q.UserNameQuery.GetUserName(), object.TextMethodToQuery(q.UserNameQuery.GetMethod()))
	case *user.SearchQuery_FirstNameQuery:
		return query.NewUserFirstNameSearchQuery(q.FirstNameQuery.GetFirstName(), object.TextMethodToQuery(q.FirstNameQuery.GetMethod()))
	case *user.SearchQuery_LastNameQuery:
		return query.NewUserLastNameSearchQuery(q.LastNameQuery.GetLastName(), object.TextMethodToQuery(q.LastNameQuery.GetMethod()))
	case *user.SearchQuery_NickNameQuery:
		return query.NewUserNickNameSearchQuery(q.NickNameQuery.GetNickName(), object.TextMethodToQuery(q.NickNameQuery.GetMethod()))
	case *user.SearchQuery_DisplayNameQuery:
		return query.NewUserDisplayNameSearchQuery(q.DisplayNameQuery.GetDisplayName(), object.TextMethodToQuery(q.DisplayNameQuery.GetMethod()))
	case *user.SearchQuery_EmailQuery:
		return query.NewUserEmailSearchQuery(q.EmailQuery.GetEmailAddress(), object.TextMethodToQuery(q.EmailQuery.GetMethod()))
	case *user.SearchQuery_PhoneQuery:
		return query.NewUserPhoneSearchQuery(q.PhoneQuery.GetNumber(), object.TextMethodToQuery(q.PhoneQuery.GetMethod()))
	case *user.SearchQuery_StateQuery:
		return query.NewUserStateSearchQuery(q.StateQuery.GetState().ToDomain())
	case *user.SearchQuery_TypeQuery:
		return query.NewUserTypeSearchQuery(q.TypeQuery.GetType().ToDomain())
	case *user.SearchQuery_LoginNameQuery:
		return query.NewUserLoginNameExistsQuery(q.LoginNameQuery.GetLoginName(), object.TextMethodToQuery(q.LoginNameQuery.GetMethod()))
	case *user.SearchQuery_OrganizationIdQuery:
		return query.NewUserResourceOwnerSearchQuery(q.OrganizationIdQuery.GetOrganizationId(), query.TextEquals)
	case *user.SearchQuery_InUserIdsQuery:
		return query.NewUserInUserIdsSearchQuery(q.InUserIdsQuery.GetUserIds())
	case *user.SearchQuery_OrQuery:
		mappedQueries, err := userQueriesToQuery(q.OrQuery.GetQueries(), level+1)
		if err != nil {
			return nil, err
		}
		return query.NewUserOrSearchQuery(mappedQueries)
	case *user.SearchQuery_AndQuery:
		mappedQueries, err := userQueriesToQuery(q.AndQuery.GetQueries(), level+1)
		if err != nil {
			return nil, err
		}
		return query.NewUserAndSearchQuery(mappedQueries)
	case *user.SearchQuery_NotQuery:
		mappedQuery, err := userQueryToQuery(q.NotQuery.GetQuery(), level+1)
		if err != nil {
			return nil, err
		}
		return query.NewUserNotSearchQuery(mappedQuery)
	case *user.SearchQuery_InUserEmailsQuery:
		return query.NewUserInUserEmailsSearchQuery(q.InUserEmailsQuery.GetUserEmails())
	case *user.SearchQuery_MetadataKeyFilter:
		return query.NewUserMetadataKeySearchQuery(q.MetadataKeyFilter.GetKey(), query.TextComparison(q.MetadataKeyFilter.GetMethod()))
	case *user.SearchQuery_MetadataValueFilter:
		return query.NewUserMetadataValueSearchQuery(q.MetadataValueFilter.GetValue(), query.BytesComparison(q.MetadataValueFilter.GetMethod()))
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-vR9nC", "List.Query.Invalid")
	}
}

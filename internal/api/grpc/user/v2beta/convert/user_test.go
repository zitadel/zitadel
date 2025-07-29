package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func TestUsersToPb(t *testing.T) {
	t.Parallel()
	users := []*query.User{
		{
			ID:                 "user1",
			Sequence:           1,
			ChangeDate:         time.Now(),
			ResourceOwner:      "owner1",
			State:              domain.UserStateActive,
			Username:           "username1",
			LoginNames:         []string{"login1"},
			PreferredLoginName: "preferred1",
			Human: &query.Human{
				FirstName:              "John",
				LastName:               "Doe",
				NickName:               "JD",
				DisplayName:            "John Doe",
				PreferredLanguage:      language.English,
				Gender:                 domain.GenderMale,
				AvatarKey:              "avatar1",
				Email:                  "john@example.com",
				IsEmailVerified:        true,
				Phone:                  "123456789",
				IsPhoneVerified:        false,
				PasswordChangeRequired: true,
				PasswordChanged:        time.Now(),
			},
		},
	}
	pbUsers := UsersToPb(users, "prefix")
	require.Len(t, pbUsers, 1)
	assert.Equal(t, "user1", pbUsers[0].UserId)
	assert.Equal(t, user.UserState_USER_STATE_ACTIVE, pbUsers[0].State)
	assert.NotNil(t, pbUsers[0].Type)
}

func TestUserToPb_Human(t *testing.T) {
	t.Parallel()
	now := time.Now()
	u := &query.User{
		ID:                 "id",
		Sequence:           2,
		ChangeDate:         now,
		ResourceOwner:      "owner",
		State:              domain.UserStateInactive,
		Username:           "uname",
		LoginNames:         []string{"ln"},
		PreferredLoginName: "pln",
		Human: &query.Human{
			FirstName:              "Jane",
			LastName:               "Smith",
			NickName:               "JS",
			DisplayName:            "Jane Smith",
			PreferredLanguage:      language.German,
			Gender:                 domain.GenderFemale,
			AvatarKey:              "avatar2",
			Email:                  "jane@example.com",
			IsEmailVerified:        false,
			Phone:                  "987654321",
			IsPhoneVerified:        true,
			PasswordChangeRequired: false,
			PasswordChanged:        now,
		},
	}
	pb := UserToPb(u, "prefix")
	assert.Equal(t, "id", pb.UserId)
	assert.Equal(t, user.UserState_USER_STATE_INACTIVE, pb.State)
	assert.NotNil(t, pb.Type)
	human, ok := pb.Type.(*user.User_Human)
	require.True(t, ok)
	assert.Equal(t, "Jane", human.Human.Profile.GivenName)
	assert.Equal(t, "Smith", human.Human.Profile.FamilyName)
	assert.Equal(t, "avatar2", u.Human.AvatarKey)
	assert.Equal(t, "jane@example.com", human.Human.Email.Email)
	assert.Equal(t, false, human.Human.Email.IsVerified)
	assert.Equal(t, "987654321", human.Human.Phone.Phone)
	assert.Equal(t, true, human.Human.Phone.IsVerified)
	assert.Equal(t, false, human.Human.PasswordChangeRequired)
	assert.Equal(t, timestamppb.New(now), human.Human.PasswordChanged)
}

func TestUserToPb_Machine(t *testing.T) {
	t.Parallel()
	u := &query.User{
		ID:                 "id2",
		Sequence:           3,
		ChangeDate:         time.Now(),
		ResourceOwner:      "owner2",
		State:              domain.UserStateDeleted,
		Username:           "uname2",
		LoginNames:         []string{"ln2"},
		PreferredLoginName: "pln2",
		Machine: &query.Machine{
			Name:            "machine1",
			Description:     "desc",
			EncodedSecret:   "secret",
			AccessTokenType: domain.OIDCTokenTypeJWT,
		},
	}
	pb := UserToPb(u, "prefix")
	assert.Equal(t, "id2", pb.UserId)
	assert.Equal(t, user.UserState_USER_STATE_DELETED, pb.State)
	assert.NotNil(t, pb.Type)
	machine, ok := pb.Type.(*user.User_Machine)
	require.True(t, ok)
	assert.Equal(t, "machine1", machine.Machine.Name)
	assert.Equal(t, "desc", machine.Machine.Description)
	assert.True(t, machine.Machine.HasSecret)
	assert.Equal(t, user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT, machine.Machine.AccessTokenType)
}

func TestUserTypeToPb_Nil(t *testing.T) {
	t.Parallel()
	u := &query.User{}
	assert.Nil(t, userTypeToPb(u, "prefix"))
}

func TestHumanToPb_PasswordChangedZero(t *testing.T) {
	t.Parallel()
	h := &query.Human{
		FirstName:              "A",
		LastName:               "B",
		NickName:               "C",
		DisplayName:            "D",
		PreferredLanguage:      language.French,
		Gender:                 domain.GenderDiverse,
		AvatarKey:              "avatar",
		Email:                  "a@b.com",
		IsEmailVerified:        true,
		Phone:                  "123",
		IsPhoneVerified:        false,
		PasswordChangeRequired: true,
		PasswordChanged:        time.Time{},
	}
	pb := humanToPb(h, "prefix", "owner")
	assert.Nil(t, pb.PasswordChanged)
}

func TestMachineToPb(t *testing.T) {
	t.Parallel()
	m := &query.Machine{
		Name:            "mach",
		Description:     "desc",
		EncodedSecret:   "",
		AccessTokenType: domain.OIDCTokenTypeBearer,
	}
	pb := machineToPb(m)
	assert.Equal(t, "mach", pb.Name)
	assert.Equal(t, "desc", pb.Description)
	assert.False(t, pb.HasSecret)
	assert.Equal(t, user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER, pb.AccessTokenType)
}

func TestUserStateToPb(t *testing.T) {
	t.Parallel()
	assert.Equal(t, user.UserState_USER_STATE_ACTIVE, userStateToPb(domain.UserStateActive))
	assert.Equal(t, user.UserState_USER_STATE_INACTIVE, userStateToPb(domain.UserStateInactive))
	assert.Equal(t, user.UserState_USER_STATE_DELETED, userStateToPb(domain.UserStateDeleted))
	assert.Equal(t, user.UserState_USER_STATE_INITIAL, userStateToPb(domain.UserStateInitial))
	assert.Equal(t, user.UserState_USER_STATE_LOCKED, userStateToPb(domain.UserStateLocked))
	assert.Equal(t, user.UserState_USER_STATE_UNSPECIFIED, userStateToPb(domain.UserStateUnspecified))
	assert.Equal(t, user.UserState_USER_STATE_UNSPECIFIED, userStateToPb(domain.UserStateSuspend))
	assert.Equal(t, user.UserState_USER_STATE_UNSPECIFIED, userStateToPb(999))
}

func TestGenderToPb(t *testing.T) {
	t.Parallel()
	assert.Equal(t, user.Gender_GENDER_DIVERSE, genderToPb(domain.GenderDiverse))
	assert.Equal(t, user.Gender_GENDER_FEMALE, genderToPb(domain.GenderFemale))
	assert.Equal(t, user.Gender_GENDER_MALE, genderToPb(domain.GenderMale))
	assert.Equal(t, user.Gender_GENDER_UNSPECIFIED, genderToPb(domain.GenderUnspecified))
	assert.Equal(t, user.Gender_GENDER_UNSPECIFIED, genderToPb(999))
}

func TestAccessTokenTypeToPb(t *testing.T) {
	t.Parallel()
	assert.Equal(t, user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER, accessTokenTypeToPb(domain.OIDCTokenTypeBearer))
	assert.Equal(t, user.AccessTokenType_ACCESS_TOKEN_TYPE_JWT, accessTokenTypeToPb(domain.OIDCTokenTypeJWT))
	assert.Equal(t, user.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER, accessTokenTypeToPb(999))
}

func TestUserFieldNameToSortingColumn(t *testing.T) {
	t.Parallel()
	assert.Equal(t, query.HumanEmailCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_EMAIL))
	assert.Equal(t, query.HumanFirstNameCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_FIRST_NAME))
	assert.Equal(t, query.HumanLastNameCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_LAST_NAME))
	assert.Equal(t, query.HumanDisplayNameCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_DISPLAY_NAME))
	assert.Equal(t, query.UserUsernameCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_USER_NAME))
	assert.Equal(t, query.UserStateCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_STATE))
	assert.Equal(t, query.UserTypeCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_TYPE))
	assert.Equal(t, query.HumanNickNameCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_NICK_NAME))
	assert.Equal(t, query.UserCreationDateCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_CREATION_DATE))
	assert.Equal(t, query.UserIDCol, userFieldNameToSortingColumn(user.UserFieldName_USER_FIELD_NAME_UNSPECIFIED))
	assert.Equal(t, query.UserIDCol, userFieldNameToSortingColumn(999))
}

func TestListUsersRequestToModel(t *testing.T) {
	t.Parallel()
	req := &user.ListUsersRequest{
		Query: &object.ListQuery{
			Offset: 1,
			Limit:  2,
			Asc:    true,
		},
		SortingColumn: user.UserFieldName_USER_FIELD_NAME_EMAIL,
		Queries: []*user.SearchQuery{
			{
				Query: &user.SearchQuery_UserNameQuery{
					UserNameQuery: &user.UserNameQuery{
						UserName: "test",
						Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
		},
	}
	model, err := ListUsersRequestToModel(req)
	require.NoError(t, err)
	assert.EqualValues(t, 1, model.Offset)
	assert.EqualValues(t, 2, model.Limit)
	assert.True(t, model.Asc)
	assert.Equal(t, query.HumanEmailCol, model.SortingColumn)
	require.Len(t, model.Queries, 1)
}

func TestUserQueriesToQuery_TooDeep(t *testing.T) {
	t.Parallel()
	q := &user.SearchQuery{
		Query: &user.SearchQuery_OrQuery{
			OrQuery: &user.OrQuery{
				Queries: []*user.SearchQuery{},
			},
		},
	}
	_, err := userQueryToQuery(q, 21)
	assert.Error(t, err)
}

func TestUserQueryToQuery_Invalid(t *testing.T) {
	t.Parallel()
	q := &user.SearchQuery{}
	_, err := userQueryToQuery(q, 0)
	assert.Error(t, err)
}

func TestUserQueryToQuery_AllTypes(t *testing.T) {
	t.Parallel()
	queries := []*user.SearchQuery{
		{Query: &user.SearchQuery_UserNameQuery{UserNameQuery: &user.UserNameQuery{UserName: "u", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_FirstNameQuery{FirstNameQuery: &user.FirstNameQuery{FirstName: "f", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_LastNameQuery{LastNameQuery: &user.LastNameQuery{LastName: "l", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_NickNameQuery{NickNameQuery: &user.NickNameQuery{NickName: "n", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_DisplayNameQuery{DisplayNameQuery: &user.DisplayNameQuery{DisplayName: "d", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_EmailQuery{EmailQuery: &user.EmailQuery{EmailAddress: "e", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_PhoneQuery{PhoneQuery: &user.PhoneQuery{Number: "p", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_StateQuery{StateQuery: &user.StateQuery{State: user.UserState_USER_STATE_ACTIVE}}},
		{Query: &user.SearchQuery_TypeQuery{TypeQuery: &user.TypeQuery{Type: user.Type_TYPE_HUMAN}}},
		{Query: &user.SearchQuery_LoginNameQuery{LoginNameQuery: &user.LoginNameQuery{LoginName: "ln", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		{Query: &user.SearchQuery_OrganizationIdQuery{OrganizationIdQuery: &user.OrganizationIdQuery{OrganizationId: "org"}}},
		{Query: &user.SearchQuery_InUserIdsQuery{InUserIdsQuery: &user.InUserIDQuery{UserIds: []string{"id"}}}},
		{Query: &user.SearchQuery_OrQuery{OrQuery: &user.OrQuery{Queries: []*user.SearchQuery{
			{Query: &user.SearchQuery_UserNameQuery{UserNameQuery: &user.UserNameQuery{UserName: "u", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
			{Query: &user.SearchQuery_DisplayNameQuery{DisplayNameQuery: &user.DisplayNameQuery{DisplayName: "d", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		}}}},
		{Query: &user.SearchQuery_AndQuery{AndQuery: &user.AndQuery{Queries: []*user.SearchQuery{
			{Query: &user.SearchQuery_UserNameQuery{UserNameQuery: &user.UserNameQuery{UserName: "u", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
			{Query: &user.SearchQuery_DisplayNameQuery{DisplayNameQuery: &user.DisplayNameQuery{DisplayName: "d", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}},
		}}}},
		{Query: &user.SearchQuery_NotQuery{NotQuery: &user.NotQuery{Query: &user.SearchQuery{Query: &user.SearchQuery_UserNameQuery{UserNameQuery: &user.UserNameQuery{UserName: "u", Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS}}}}}},
		{Query: &user.SearchQuery_InUserEmailsQuery{InUserEmailsQuery: &user.InUserEmailsQuery{UserEmails: []string{"e"}}}},
	}
	for _, q := range queries {
		_, err := userQueryToQuery(q, 0)
		assert.NoError(t, err)
	}
}

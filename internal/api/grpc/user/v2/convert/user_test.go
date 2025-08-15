package convert

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	metadata "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func Test_userQueryToQuery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		input   *user.SearchQuery
		level   uint8
		wantErr bool
	}{
		{
			name: "UserNameQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_UserNameQuery{
					UserNameQuery: &user.UserNameQuery{
						UserName: "testuser",
						Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "FirstNameQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_FirstNameQuery{
					FirstNameQuery: &user.FirstNameQuery{
						FirstName: "John",
						Method:    object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "LastNameQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_LastNameQuery{
					LastNameQuery: &user.LastNameQuery{
						LastName: "Doe",
						Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "NickNameQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_NickNameQuery{
					NickNameQuery: &user.NickNameQuery{
						NickName: "JD",
						Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "DisplayNameQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_DisplayNameQuery{
					DisplayNameQuery: &user.DisplayNameQuery{
						DisplayName: "John Doe",
						Method:      object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "EmailQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_EmailQuery{
					EmailQuery: &user.EmailQuery{
						EmailAddress: "john@doe.com",
						Method:       object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "PhoneQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_PhoneQuery{
					PhoneQuery: &user.PhoneQuery{
						Number: "123456789",
						Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "StateQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_StateQuery{
					StateQuery: &user.StateQuery{
						State: user.UserState_USER_STATE_ACTIVE,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "TypeQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_TypeQuery{
					TypeQuery: &user.TypeQuery{
						Type: user.Type_TYPE_HUMAN,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "LoginNameQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_LoginNameQuery{
					LoginNameQuery: &user.LoginNameQuery{
						LoginName: "login",
						Method:    object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "OrganizationIdQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_OrganizationIdQuery{
					OrganizationIdQuery: &user.OrganizationIdQuery{
						OrganizationId: "org1",
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "InUserIdsQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_InUserIdsQuery{
					InUserIdsQuery: &user.InUserIDQuery{
						UserIds: []string{"id1", "id2"},
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "OrQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_OrQuery{
					OrQuery: &user.OrQuery{
						Queries: []*user.SearchQuery{
							{
								Query: &user.SearchQuery_UserNameQuery{
									UserNameQuery: &user.UserNameQuery{
										UserName: "testuser",
										Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
									},
								},
							},
						},
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "AndQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_AndQuery{
					AndQuery: &user.AndQuery{
						Queries: []*user.SearchQuery{
							{
								Query: &user.SearchQuery_UserNameQuery{
									UserNameQuery: &user.UserNameQuery{
										UserName: "testuser",
										Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
									},
								},
							},
						},
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "NotQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_NotQuery{
					NotQuery: &user.NotQuery{
						Query: &user.SearchQuery{
							Query: &user.SearchQuery_UserNameQuery{
								UserNameQuery: &user.UserNameQuery{
									UserName: "testuser",
									Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
								},
							},
						},
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "InUserEmailsQuery",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_InUserEmailsQuery{
					InUserEmailsQuery: &user.InUserEmailsQuery{
						UserEmails: []string{"john@doe.com", "jane@doe.com"},
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "MetadataKeyFilter",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_MetadataKeyFilter{
					MetadataKeyFilter: &metadata.MetadataKeyFilter{
						Key:    "key1",
						Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "MetadataValueFilter",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_MetadataValueFilter{
					MetadataValueFilter: &metadata.MetadataValueFilter{
						Value:  []byte("value1"),
						Method: filter.ByteFilterMethod_BYTE_FILTER_METHOD_EQUALS,
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "TooManyNestingLevels",
			input: &user.SearchQuery{
				Query: &user.SearchQuery_UserNameQuery{
					UserNameQuery: &user.UserNameQuery{
						UserName: "testuser",
						Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
			level:   21,
			wantErr: true,
		},
		{
			name: "InvalidQueryType",
			input: &user.SearchQuery{
				Query: nil,
			},
			level:   0,
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := userQueryToQuery(tc.input, tc.level)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_userQueriesToQuery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		queries []*user.SearchQuery
		level   uint8
		wantErr bool
	}{
		{
			name: "single valid query",
			queries: []*user.SearchQuery{
				{
					Query: &user.SearchQuery_UserNameQuery{
						UserNameQuery: &user.UserNameQuery{
							UserName: "testuser",
							Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name: "multiple valid queries",
			queries: []*user.SearchQuery{
				{
					Query: &user.SearchQuery_UserNameQuery{
						UserNameQuery: &user.UserNameQuery{
							UserName: "user1",
							Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
				{
					Query: &user.SearchQuery_EmailQuery{
						EmailQuery: &user.EmailQuery{
							EmailAddress: "user1@example.com",
							Method:       object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
			},
			level:   0,
			wantErr: false,
		},
		{
			name:    "empty queries slice",
			queries: []*user.SearchQuery{},
			level:   0,
			wantErr: false,
		},
		{
			name: "query with error (too many nesting levels)",
			queries: []*user.SearchQuery{
				{
					Query: &user.SearchQuery_UserNameQuery{
						UserNameQuery: &user.UserNameQuery{
							UserName: "testuser",
							Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
			},
			level:   21,
			wantErr: true,
		},
		{
			name: "query with invalid type",
			queries: []*user.SearchQuery{
				{
					Query: nil,
				},
			},
			level:   0,
			wantErr: true,
		},
		{
			name: "mixed valid and invalid queries",
			queries: []*user.SearchQuery{
				{
					Query: &user.SearchQuery_UserNameQuery{
						UserNameQuery: &user.UserNameQuery{
							UserName: "testuser",
							Method:   object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
						},
					},
				},
				{
					Query: nil,
				},
			},
			level:   0,
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := userQueriesToQuery(tc.queries, tc.level)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, len(tc.queries), len(got))
			}
		})
	}
}

func Test_userFieldNameToSortingColumn(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		input    user.UserFieldName
		expected query.Column
	}

	tt := []testCase{
		{
			name:     "EMAIL",
			input:    user.UserFieldName_USER_FIELD_NAME_EMAIL,
			expected: query.HumanEmailCol,
		},
		{
			name:     "FIRST_NAME",
			input:    user.UserFieldName_USER_FIELD_NAME_FIRST_NAME,
			expected: query.HumanFirstNameCol,
		},
		{
			name:     "LAST_NAME",
			input:    user.UserFieldName_USER_FIELD_NAME_LAST_NAME,
			expected: query.HumanLastNameCol,
		},
		{
			name:     "DISPLAY_NAME",
			input:    user.UserFieldName_USER_FIELD_NAME_DISPLAY_NAME,
			expected: query.HumanDisplayNameCol,
		},
		{
			name:     "USER_NAME",
			input:    user.UserFieldName_USER_FIELD_NAME_USER_NAME,
			expected: query.UserUsernameCol,
		},
		{
			name:     "STATE",
			input:    user.UserFieldName_USER_FIELD_NAME_STATE,
			expected: query.UserStateCol,
		},
		{
			name:     "TYPE",
			input:    user.UserFieldName_USER_FIELD_NAME_TYPE,
			expected: query.UserTypeCol,
		},
		{
			name:     "NICK_NAME",
			input:    user.UserFieldName_USER_FIELD_NAME_NICK_NAME,
			expected: query.HumanNickNameCol,
		},
		{
			name:     "CREATION_DATE",
			input:    user.UserFieldName_USER_FIELD_NAME_CREATION_DATE,
			expected: query.UserCreationDateCol,
		},
		{
			name:     "UNSPECIFIED",
			input:    user.UserFieldName_USER_FIELD_NAME_UNSPECIFIED,
			expected: query.UserIDCol,
		},
		{
			name:     "Unknown value (default)",
			input:    user.UserFieldName(999),
			expected: query.UserIDCol,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := userFieldNameToSortingColumn(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_userStateToPb(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		input    domain.UserState
		expected user.UserState
	}

	tt := []testCase{
		{
			name:     "Active",
			input:    domain.UserStateActive,
			expected: user.UserState_USER_STATE_ACTIVE,
		},
		{
			name:     "Inactive",
			input:    domain.UserStateInactive,
			expected: user.UserState_USER_STATE_INACTIVE,
		},
		{
			name:     "Deleted",
			input:    domain.UserStateDeleted,
			expected: user.UserState_USER_STATE_DELETED,
		},
		{
			name:     "Initial",
			input:    domain.UserStateInitial,
			expected: user.UserState_USER_STATE_INITIAL,
		},
		{
			name:     "Locked",
			input:    domain.UserStateLocked,
			expected: user.UserState_USER_STATE_LOCKED,
		},
		{
			name:     "Unspecified",
			input:    domain.UserStateUnspecified,
			expected: user.UserState_USER_STATE_UNSPECIFIED,
		},
		{
			name:     "Suspend",
			input:    domain.UserStateSuspend,
			expected: user.UserState_USER_STATE_UNSPECIFIED,
		},
		{
			name:     "Unknown value (default)",
			input:    domain.UserState(999),
			expected: user.UserState_USER_STATE_UNSPECIFIED,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := userStateToPb(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_userTypeToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		testName         string
		inputQueryUser   *query.User
		expectedUserType user.UserType
	}{
		{
			testName:         "when query user is nil should return nil",
			expectedUserType: nil,
		},
		{
			testName:       "when query user is human should return usertype human",
			inputQueryUser: &query.User{Human: &query.Human{}},
			expectedUserType: &user.User_Human{
				Human: &user.HumanUser{
					Profile: &user.HumanProfile{
						NickName:          new(string),
						DisplayName:       new(string),
						PreferredLanguage: gu.Ptr("und"),
						Gender:            gu.Ptr(user.Gender_GENDER_UNSPECIFIED),
					},
					Email: &user.HumanEmail{},
					Phone: &user.HumanPhone{},
				},
			},
		},
		{
			testName:       "when query user is machine should return usertype machine",
			inputQueryUser: &query.User{Machine: &query.Machine{}},
			expectedUserType: &user.User_Machine{
				Machine: &user.MachineUser{},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			res := userTypeToPb(tc.inputQueryUser, "asset")

			assert.Equal(t, tc.expectedUserType, res)
		})
	}
}

func Test_UserToPb(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	type args struct {
		userQ       *query.User
		assetPrefix string
	}
	tt := []struct {
		name string
		args args
		want *user.User
	}{
		{
			name: "nil userQ returns nil",
			args: args{
				userQ:       nil,
				assetPrefix: "prefix",
			},
			want: nil,
		},
		{
			name: "basic user with human type",
			args: args{
				userQ: &query.User{
					ID:                 "user-id",
					Sequence:           1,
					ChangeDate:         now,
					ResourceOwner:      "owner",
					CreationDate:       now,
					State:              domain.UserStateActive,
					Username:           "username",
					LoginNames:         []string{"login1", "login2"},
					PreferredLoginName: "preferred",
					Human:              &query.Human{},
				},
				assetPrefix: "prefix",
			},
			want: &user.User{
				UserId: "user-id",
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.New(now),
					ResourceOwner: "owner",
					CreationDate:  timestamppb.New(now),
				},
				State:              user.UserState_USER_STATE_ACTIVE,
				Username:           "username",
				LoginNames:         []string{"login1", "login2"},
				PreferredLoginName: "preferred",
				Type: &user.User_Human{
					Human: humanToPb(&query.Human{}, "prefix", "owner"),
				},
			},
		},
		{
			name: "basic user with machine type",
			args: args{
				userQ: &query.User{
					ID:                 "user-id2",
					Sequence:           2,
					ChangeDate:         now,
					ResourceOwner:      "owner2",
					CreationDate:       now,
					State:              domain.UserStateInactive,
					Username:           "machineuser",
					LoginNames:         []string{"machine1"},
					PreferredLoginName: "machinepreferred",
					Machine:            &query.Machine{},
				},
				assetPrefix: "prefix2",
			},
			want: &user.User{
				UserId: "user-id2",
				Details: &object.Details{
					Sequence:      2,
					ChangeDate:    timestamppb.New(now),
					ResourceOwner: "owner2",
					CreationDate:  timestamppb.New(now),
				},
				State:              user.UserState_USER_STATE_INACTIVE,
				Username:           "machineuser",
				LoginNames:         []string{"machine1"},
				PreferredLoginName: "machinepreferred",
				Type: &user.User_Machine{
					Machine: machineToPb(&query.Machine{}),
				},
			},
		},
		{
			name: "user with no type returns nil Type",
			args: args{
				userQ: &query.User{
					ID:                 "user-id3",
					Sequence:           3,
					ChangeDate:         now,
					ResourceOwner:      "owner3",
					CreationDate:       now,
					State:              domain.UserStateDeleted,
					Username:           "notypeuser",
					LoginNames:         []string{"notype1"},
					PreferredLoginName: "notypepreferred",
				},
				assetPrefix: "prefix3",
			},
			want: &user.User{
				UserId: "user-id3",
				Details: &object.Details{
					Sequence:      3,
					ChangeDate:    timestamppb.New(now),
					ResourceOwner: "owner3",
					CreationDate:  timestamppb.New(now),
				},
				State:              user.UserState_USER_STATE_DELETED,
				Username:           "notypeuser",
				LoginNames:         []string{"notype1"},
				PreferredLoginName: "notypepreferred",
				Type:               nil,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := UserToPb(tc.args.userQ, tc.args.assetPrefix)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_UsersToPb(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()

	type args struct {
		users       []*query.User
		assetPrefix string
	}
	tt := []struct {
		name string
		args args
		want []*user.User
	}{
		{
			name: "nil slice returns nils",
			args: args{
				users:       nil,
				assetPrefix: "prefix",
			},
			want: []*user.User{},
		},
		{
			name: "empty slice returns empty slice",
			args: args{
				users:       []*query.User{},
				assetPrefix: "prefix",
			},
			want: []*user.User{},
		},
		{
			name: "slice with nil user returns slice with nil",
			args: args{
				users:       []*query.User{nil},
				assetPrefix: "prefix",
			},
			want: []*user.User{nil},
		},
		{
			name: "slice with one human user",
			args: args{
				users: []*query.User{
					{
						ID:                 "user-id",
						Sequence:           1,
						ChangeDate:         now,
						ResourceOwner:      "owner",
						CreationDate:       now,
						State:              domain.UserStateActive,
						Username:           "username",
						LoginNames:         []string{"login1", "login2"},
						PreferredLoginName: "preferred",
						Human:              &query.Human{},
					},
				},
				assetPrefix: "prefix",
			},
			want: []*user.User{
				{
					UserId: "user-id",
					Details: &object.Details{
						Sequence:      1,
						ChangeDate:    timestamppb.New(now),
						ResourceOwner: "owner",
						CreationDate:  timestamppb.New(now),
					},
					State:              user.UserState_USER_STATE_ACTIVE,
					Username:           "username",
					LoginNames:         []string{"login1", "login2"},
					PreferredLoginName: "preferred",
					Type: &user.User_Human{
						Human: humanToPb(&query.Human{}, "prefix", "owner"),
					},
				},
			},
		},
		{
			name: "slice with one machine user",
			args: args{
				users: []*query.User{
					{
						ID:                 "user-id2",
						Sequence:           2,
						ChangeDate:         now,
						ResourceOwner:      "owner2",
						CreationDate:       now,
						State:              domain.UserStateInactive,
						Username:           "machineuser",
						LoginNames:         []string{"machine1"},
						PreferredLoginName: "machinepreferred",
						Machine:            &query.Machine{},
					},
				},
				assetPrefix: "prefix2",
			},
			want: []*user.User{
				{
					UserId: "user-id2",
					Details: &object.Details{
						Sequence:      2,
						ChangeDate:    timestamppb.New(now),
						ResourceOwner: "owner2",
						CreationDate:  timestamppb.New(now),
					},
					State:              user.UserState_USER_STATE_INACTIVE,
					Username:           "machineuser",
					LoginNames:         []string{"machine1"},
					PreferredLoginName: "machinepreferred",
					Type: &user.User_Machine{
						Machine: machineToPb(&query.Machine{}),
					},
				},
			},
		},
		{
			name: "slice with mixed users",
			args: args{
				users: []*query.User{
					nil,
					{
						ID:                 "user-id3",
						Sequence:           3,
						ChangeDate:         now,
						ResourceOwner:      "owner3",
						CreationDate:       now,
						State:              domain.UserStateDeleted,
						Username:           "notypeuser",
						LoginNames:         []string{"notype1"},
						PreferredLoginName: "notypepreferred",
					},
				},
				assetPrefix: "prefix3",
			},
			want: []*user.User{
				nil,
				{
					UserId: "user-id3",
					Details: &object.Details{
						Sequence:      3,
						ChangeDate:    timestamppb.New(now),
						ResourceOwner: "owner3",
						CreationDate:  timestamppb.New(now),
					},
					State:              user.UserState_USER_STATE_DELETED,
					Username:           "notypeuser",
					LoginNames:         []string{"notype1"},
					PreferredLoginName: "notypepreferred",
					Type:               nil,
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := UsersToPb(tc.args.users, tc.args.assetPrefix)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_ListUsersRequestToModel(t *testing.T) {
	t.Parallel()

	firstNameQuery, err := query.NewUserFirstNameSearchQuery("first name", query.TextEquals)
	require.Nil(t, err)

	type args struct {
		req *user.ListUsersRequest
	}
	tt := []struct {
		name    string
		args    args
		want    *query.UserSearchQueries
		wantErr bool
	}{
		{
			name: "valid request with one query",
			args: args{
				req: &user.ListUsersRequest{
					Query: &object.ListQuery{
						Offset: 10,
						Limit:  5,
						Asc:    true,
					},
					SortingColumn: user.UserFieldName_USER_FIELD_NAME_EMAIL,
					Queries: []*user.SearchQuery{
						{
							Query: &user.SearchQuery_FirstNameQuery{
								FirstNameQuery: &user.FirstNameQuery{
									FirstName: "first name",
									Method:    object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
								},
							},
						},
					},
				},
			},
			want: &query.UserSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        10,
					Limit:         5,
					Asc:           true,
					SortingColumn: query.HumanEmailCol,
				},
				Queries: []query.SearchQuery{firstNameQuery},
			},
			wantErr: false,
		},
		{
			name: "valid request with no queries",
			args: args{
				req: &user.ListUsersRequest{
					Query:         &object.ListQuery{Offset: 0, Limit: 0, Asc: false},
					SortingColumn: user.UserFieldName_USER_FIELD_NAME_UNSPECIFIED,
					Queries:       []*user.SearchQuery{},
				},
			},
			want: &query.UserSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         0,
					Asc:           false,
					SortingColumn: query.UserIDCol,
				},
				Queries: []query.SearchQuery{},
			},
			wantErr: false,
		},
		{
			name: "invalid query type returns error",
			args: args{
				req: &user.ListUsersRequest{
					Query:         &object.ListQuery{Offset: 0, Limit: 0, Asc: false},
					SortingColumn: user.UserFieldName_USER_FIELD_NAME_EMAIL,
					Queries: []*user.SearchQuery{
						{Query: nil},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil request returns zero values",
			args: args{
				req: nil,
			},
			want: &query.UserSearchQueries{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         0,
					Asc:           false,
					SortingColumn: query.UserIDCol,
				},
				Queries: []query.SearchQuery{},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := ListUsersRequestToModel(tc.args.req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/query"
	filter "github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	metadata "github.com/zitadel/zitadel/pkg/grpc/metadata/v2"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func Test_usersByMetadataSorting(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name     string
		input    user.UsersByMetadataSorting
		expected query.Column
	}{
		{"DisplayName", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_DISPLAY_NAME, query.HumanDisplayNameCol},
		{"Email", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_EMAIL, query.HumanEmailCol},
		{"FirstName", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_FIRST_NAME, query.HumanFirstNameCol},
		{"LastName", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_LAST_NAME, query.HumanLastNameCol},
		{"MetadataKey", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_METADATA_KEY, query.UserMetadataKeyCol},
		{"MetadataValue", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_METADATA_VALUE, query.UserMetadataValueCol},
		{"NickName", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_NICK_NAME, query.HumanNickNameCol},
		{"State", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_STATE, query.UserStateCol},
		{"Type", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_TYPE, query.UserTypeCol},
		{"UserName", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_USER_NAME, query.UserUsernameCol},
		{"UserID", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_USER_ID, query.UserIDCol},
		{"Unspecified", user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_UNSPECIFIED, query.UserIDCol},
		{"Default", user.UsersByMetadataSorting(999), query.UserIDCol},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := usersByMetadataSorting(tc.input)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_userByMetadataQuery(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name    string
		input   *metadata.UserByMetadataSearchFilter
		wantErr bool
	}{
		{
			name: "KeyFilter",
			input: &metadata.UserByMetadataSearchFilter{
				Filter: &metadata.UserByMetadataSearchFilter_KeyFilter{
					KeyFilter: &metadata.MetadataKeyFilter{
						Key:    "foo",
						Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ValueFilter",
			input: &metadata.UserByMetadataSearchFilter{
				Filter: &metadata.UserByMetadataSearchFilter_ValueFilter{
					ValueFilter: &metadata.MetadataValueFilter{
						Value:  []byte("bar"),
						Method: filter.ByteFilterMethod_BYTE_FILTER_METHOD_EQUALS,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "AndFilter",
			input: &metadata.UserByMetadataSearchFilter{
				Filter: &metadata.UserByMetadataSearchFilter_AndFilter{
					AndFilter: &metadata.MetadataAndFilter{
						Queries: []*metadata.UserByMetadataSearchFilter{
							{
								Filter: &metadata.UserByMetadataSearchFilter_KeyFilter{
									KeyFilter: &metadata.MetadataKeyFilter{
										Key:    "foo",
										Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "OrFilter",
			input: &metadata.UserByMetadataSearchFilter{
				Filter: &metadata.UserByMetadataSearchFilter_OrFilter{
					OrFilter: &metadata.MetadataOrFilter{
						Queries: []*metadata.UserByMetadataSearchFilter{
							{
								Filter: &metadata.UserByMetadataSearchFilter_ValueFilter{
									ValueFilter: &metadata.MetadataValueFilter{
										Value:  []byte("baz"),
										Method: filter.ByteFilterMethod_BYTE_FILTER_METHOD_EQUALS,
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "NotFilter",
			input: &metadata.UserByMetadataSearchFilter{
				Filter: &metadata.UserByMetadataSearchFilter_NotFilter{
					NotFilter: &metadata.MetadataNotFilter{
						Query: &metadata.UserByMetadataSearchFilter{
							Filter: &metadata.UserByMetadataSearchFilter_KeyFilter{
								KeyFilter: &metadata.MetadataKeyFilter{
									Key:    "not",
									Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "TooManyNestingLevels",
			input: &metadata.UserByMetadataSearchFilter{
				Filter: &metadata.UserByMetadataSearchFilter_AndFilter{
					AndFilter: &metadata.MetadataAndFilter{
						Queries: []*metadata.UserByMetadataSearchFilter{},
					},
				},
			},
			wantErr: true,
		},
		{
			name:    "InvalidFilter",
			input:   &metadata.UserByMetadataSearchFilter{},
			wantErr: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			nesting := uint(0)
			if tc.name == "TooManyNestingLevels" {
				nesting = 21
			}
			got, err := userByMetadataQuery(tc.input, nesting)
			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, got)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func Test_usersByMetadataQueries(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		queries, err := usersByMetadataQueries([]*metadata.UserByMetadataSearchFilter{}, 0)
		require.NoError(t, err)
		assert.Len(t, queries, 0)
	})

	t.Run("single valid", func(t *testing.T) {
		t.Parallel()
		queries, err := usersByMetadataQueries([]*metadata.UserByMetadataSearchFilter{
			{
				Filter: &metadata.UserByMetadataSearchFilter_KeyFilter{
					KeyFilter: &metadata.MetadataKeyFilter{
						Key:    "foo",
						Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
					},
				},
			},
		}, 0)
		require.NoError(t, err)
		assert.Len(t, queries, 1)
	})

	t.Run("invalid filter", func(t *testing.T) {
		t.Parallel()
		queries, err := usersByMetadataQueries([]*metadata.UserByMetadataSearchFilter{
			{},
		}, 0)
		require.Error(t, err)
		assert.Nil(t, queries)
	})
}

func Test_ListUsersByMetadataRequestToModel(t *testing.T) {
	t.Parallel()

	sysDefaults := systemdefaults.SystemDefaults{
		MaxQueryLimit: 100,
	}

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		req := &user.ListUsersByMetadataRequest{
			Pagination: &filter.PaginationRequest{
				Limit:  10,
				Offset: 5,
				Asc:    true,
			},
			SortingColumn: user.UsersByMetadataSorting_USERS_BY_METADATA_SORT_BY_EMAIL,
			Filters: []*metadata.UserByMetadataSearchFilter{
				{
					Filter: &metadata.UserByMetadataSearchFilter_KeyFilter{
						KeyFilter: &metadata.MetadataKeyFilter{
							Key:    "foo",
							Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_EQUALS,
						},
					},
				},
			},
		}
		model, err := ListUsersByMetadataRequestToModel(req, sysDefaults)
		require.NoError(t, err)
		assert.NotNil(t, model)
		assert.EqualValues(t, 5, model.Offset)
		assert.EqualValues(t, 10, model.Limit)
		assert.True(t, model.Asc)
		assert.Equal(t, query.HumanEmailCol, model.SortingColumn)
		assert.Len(t, model.Queries, 1)
	})

	t.Run("invalid pagination", func(t *testing.T) {
		t.Parallel()
		req := &user.ListUsersByMetadataRequest{
			Pagination: &filter.PaginationRequest{

				Limit: 200,
			},
			Filters: []*metadata.UserByMetadataSearchFilter{},
		}
		model, err := ListUsersByMetadataRequestToModel(req, sysDefaults)
		require.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("invalid filter", func(t *testing.T) {
		t.Parallel()
		req := &user.ListUsersByMetadataRequest{
			Pagination: &filter.PaginationRequest{
				Limit: 1,
			},
			Filters: []*metadata.UserByMetadataSearchFilter{
				{},
			},
		}
		model, err := ListUsersByMetadataRequestToModel(req, sysDefaults)
		require.Error(t, err)
		assert.Nil(t, model)
	})
}

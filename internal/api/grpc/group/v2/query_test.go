package group

import (
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	authorization "github.com/zitadel/zitadel/pkg/grpc/authorization/v2beta"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

func Test_ListGroupsRequestToModel(t *testing.T) {
	t.Parallel()
	groupIDsSearchQuery, err := query.NewGroupIDsSearchQuery([]string{"group1", "group2"})
	require.NoError(t, err)

	tests := []struct {
		name          string
		maxQueryLimit uint64
		req           *group_v2.ListGroupsRequest
		wantResp      *query.GroupSearchQuery
		wantErr       error
	}{
		{
			name:          "max query limit exceeded",
			maxQueryLimit: 1,
			req: &group_v2.ListGroupsRequest{
				Pagination: &filter.PaginationRequest{
					Limit: 5,
				},
				Filters: []*group_v2.GroupsSearchFilter{
					{
						Filter: &group_v2.GroupsSearchFilter_GroupIds{
							GroupIds: &filter.InIDsFilter{
								Ids: []string{"group1", "group2"},
							},
						},
					},
				},
			},
			wantErr: zerrors.ThrowInvalidArgumentf(errors.New("given: 5, allowed: 1"), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			name: "valid request, list of group IDs, ok",
			req: &group_v2.ListGroupsRequest{
				Filters: []*group_v2.GroupsSearchFilter{
					{
						Filter: &group_v2.GroupsSearchFilter_GroupIds{
							GroupIds: &filter.InIDsFilter{
								Ids: []string{"group1", "group2"},
							},
						},
					},
				},
			},
			wantResp: &query.GroupSearchQuery{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         0,
					SortingColumn: query.GroupColumnCreationDate,
				},
				Queries: []query.SearchQuery{groupIDsSearchQuery},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sysDefaults := systemdefaults.SystemDefaults{MaxQueryLimit: tt.maxQueryLimit}
			got, err := listGroupsRequestToModel(tt.req, sysDefaults)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantResp, got)
		})
	}
}

func Test_GroupSearchFiltersToQuery(t *testing.T) {
	t.Parallel()
	groupIDsSearchQuery, err := query.NewGroupIDsSearchQuery([]string{"group1", "group2"})
	require.NoError(t, err)
	groupNameSearchQuery, err := query.NewGroupNameSearchQuery("mygroup", query.TextStartsWith)
	require.NoError(t, err)
	groupOrgIDSearchQuery, err := query.NewGroupOrganizationIdSearchQuery("org1")
	require.NoError(t, err)

	tests := []struct {
		name    string
		filters []*group_v2.GroupsSearchFilter
		want    []query.SearchQuery
		wantErr error
	}{
		{
			name:    "empty",
			filters: []*group_v2.GroupsSearchFilter{},
			want:    []query.SearchQuery{},
		},
		{
			name: "all filters",
			filters: []*group_v2.GroupsSearchFilter{
				{
					Filter: &group_v2.GroupsSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{"group1", "group2"},
						},
					},
				},
				{
					Filter: &group_v2.GroupsSearchFilter_NameFilter{
						NameFilter: &group_v2.GroupNameFilter{
							Name:   "mygroup",
							Method: filter.TextFilterMethod_TEXT_FILTER_METHOD_STARTS_WITH,
						},
					},
				},
				{
					Filter: &group_v2.GroupsSearchFilter_OrganizationId{
						OrganizationId: &filter.IDFilter{
							Id: "org1",
						},
					},
				},
			},
			want: []query.SearchQuery{
				groupIDsSearchQuery,
				groupNameSearchQuery,
				groupOrgIDSearchQuery,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := groupSearchFiltersToQuery(tt.filters)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GroupFieldNameToSortingColumn(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		field *group_v2.FieldName
		want  query.Column
	}{
		{
			name:  "nil",
			field: nil,
			want:  query.GroupColumnCreationDate,
		},
		{
			name:  "creation date",
			field: gu.Ptr(group_v2.FieldName_FIELD_NAME_CREATION_DATE),
			want:  query.GroupColumnCreationDate,
		},
		{
			name:  "unspecified",
			field: gu.Ptr(group_v2.FieldName_FIELD_NAME_UNSPECIFIED),
			want:  query.GroupColumnCreationDate,
		},
		{
			name:  "id",
			field: gu.Ptr(group_v2.FieldName_FIELD_NAME_ID),
			want:  query.GroupColumnID,
		},
		{
			name:  "name",
			field: gu.Ptr(group_v2.FieldName_FIELD_NAME_NAME),
			want:  query.GroupColumnName,
		},
		{
			name:  "change date",
			field: gu.Ptr(group_v2.FieldName_FIELD_NAME_CHANGE_DATE),
			want:  query.GroupColumnChangeDate,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := groupFieldNameToSortingColumn(tt.field)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GroupsToPb(t *testing.T) {
	t.Parallel()
	timeNow := time.Now().UTC()
	tests := []struct {
		name   string
		groups []*query.Group
		want   []*group_v2.Group
	}{
		{
			name:   "empty",
			groups: []*query.Group{},
			want:   []*group_v2.Group{},
		},
		{
			name: "with groups, ok",
			groups: []*query.Group{
				{
					ID:            "group1",
					Name:          "mygroup",
					Description:   "my first group",
					CreationDate:  timeNow,
					ChangeDate:    timeNow,
					ResourceOwner: "org1",
					InstanceID:    "instance1",
					State:         domain.GroupStateActive,
					Sequence:      1,
				},
				{
					ID:            "group2",
					Name:          "mygroup2",
					Description:   "my second group",
					CreationDate:  timeNow,
					ChangeDate:    timeNow,
					ResourceOwner: "org1",
					InstanceID:    "instance1",
					State:         domain.GroupStateActive,
					Sequence:      1,
				},
			},
			want: []*group_v2.Group{
				{
					Id:             "group1",
					Name:           "mygroup",
					Description:    "my first group",
					OrganizationId: "org1",
					ChangeDate:     timestamppb.New(timeNow),
					CreationDate:   timestamppb.New(timeNow),
				},
				{
					Id:             "group2",
					Name:           "mygroup2",
					Description:    "my second group",
					OrganizationId: "org1",
					ChangeDate:     timestamppb.New(timeNow),
					CreationDate:   timestamppb.New(timeNow),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := groupsToPb(tt.groups)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_ListGroupUsersRequestToModel(t *testing.T) {
	t.Parallel()
	groupIDsSearchQuery, err := query.NewGroupUsersGroupIDsSearchQuery([]string{"group1", "group2"})
	require.NoError(t, err)

	userIDsSearchQuery, err := query.NewGroupUsersUserIDsSearchQuery([]string{"user1", "user2"})
	require.NoError(t, err)

	tests := []struct {
		name          string
		maxQueryLimit uint64
		req           *group_v2.ListGroupUsersRequest
		wantResp      *query.GroupUsersSearchQuery
		wantErr       error
	}{
		{
			name:          "max query limit exceeded",
			maxQueryLimit: 1,
			req: &group_v2.ListGroupUsersRequest{
				Pagination: &filter.PaginationRequest{
					Limit: 5,
				},
				Filters: []*group_v2.GroupUsersSearchFilter{
					{
						Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
							GroupIds: &filter.InIDsFilter{
								Ids: []string{"group1", "group2"},
							},
						},
					},
				},
			},
			wantErr: zerrors.ThrowInvalidArgumentf(errors.New("given: 5, allowed: 1"), "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
		},
		{
			name: "valid request, list of group IDs, ok",
			req: &group_v2.ListGroupUsersRequest{
				Filters: []*group_v2.GroupUsersSearchFilter{
					{
						Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
							GroupIds: &filter.InIDsFilter{
								Ids: []string{"group1", "group2"},
							},
						},
					},
				},
			},
			wantResp: &query.GroupUsersSearchQuery{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         0,
					SortingColumn: query.GroupUsersColumnCreationDate,
				},
				Queries: []query.SearchQuery{groupIDsSearchQuery},
			},
		},
		{
			name: "valid request, list of user IDs, ok",
			req: &group_v2.ListGroupUsersRequest{
				Filters: []*group_v2.GroupUsersSearchFilter{
					{
						Filter: &group_v2.GroupUsersSearchFilter_UserIds{
							UserIds: &filter.InIDsFilter{
								Ids: []string{"user1", "user2"},
							},
						},
					},
				},
			},
			wantResp: &query.GroupUsersSearchQuery{
				SearchRequest: query.SearchRequest{
					Offset:        0,
					Limit:         0,
					SortingColumn: query.GroupUsersColumnCreationDate,
				},
				Queries: []query.SearchQuery{userIDsSearchQuery},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sysDefaults := systemdefaults.SystemDefaults{MaxQueryLimit: tt.maxQueryLimit}
			got, err := listGroupUsersRequestToModel(tt.req, sysDefaults)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantResp, got)
		})
	}
}

func Test_GroupUsersSearchFiltersToQuery(t *testing.T) {
	t.Parallel()
	groupIDsSearchQuery, err := query.NewGroupUsersGroupIDsSearchQuery([]string{"group1", "group2"})
	require.NoError(t, err)

	userIDsSearchQuery, err := query.NewGroupUsersUserIDsSearchQuery([]string{"user1", "user2"})
	require.NoError(t, err)

	tests := []struct {
		name    string
		filters []*group_v2.GroupUsersSearchFilter
		want    []query.SearchQuery
		wantErr error
	}{
		{
			name:    "empty",
			filters: []*group_v2.GroupUsersSearchFilter{},
			want:    []query.SearchQuery{},
		},
		{
			name: "all filters",
			filters: []*group_v2.GroupUsersSearchFilter{
				{
					Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
						GroupIds: &filter.InIDsFilter{
							Ids: []string{"group1", "group2"},
						},
					},
				},
				{
					Filter: &group_v2.GroupUsersSearchFilter_UserIds{
						UserIds: &filter.InIDsFilter{
							Ids: []string{"user1", "user2"},
						},
					},
				},
			},
			want: []query.SearchQuery{
				groupIDsSearchQuery,
				userIDsSearchQuery,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := groupUsersSearchFiltersToQuery(tt.filters)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GroupUsersFieldNameToSortingColumn(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		field *group_v2.GroupUserFieldName
		want  query.Column
	}{
		{
			name:  "nil",
			field: nil,
			want:  query.GroupUsersColumnCreationDate,
		},
		{
			name:  "creation date",
			field: gu.Ptr(group_v2.GroupUserFieldName_GROUP_USER_FIELD_NAME_CREATION_DATE),
			want:  query.GroupUsersColumnCreationDate,
		},
		{
			name:  "unspecified",
			field: gu.Ptr(group_v2.GroupUserFieldName_GROUP_USER_FIELD_NAME_CREATION_DATE),
			want:  query.GroupUsersColumnCreationDate,
		},
		{
			name:  "group id",
			field: gu.Ptr(group_v2.GroupUserFieldName_GROUP_USER_FIELD_NAME_GROUP_ID),
			want:  query.GroupUsersColumnGroupID,
		},
		{
			name:  "user id",
			field: gu.Ptr(group_v2.GroupUserFieldName_GROUP_USER_FIELD_NAME_USER_ID),
			want:  query.GroupUsersColumnUserID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := groupUsersFieldNameToSortingColumn(tt.field)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GroupUsersToPb(t *testing.T) {
	t.Parallel()
	timeNow := time.Now().UTC()
	tests := []struct {
		name       string
		groupUsers []*query.GroupUser
		want       []*group_v2.GroupUser
	}{
		{
			name:       "empty",
			groupUsers: []*query.GroupUser{},
			want:       []*group_v2.GroupUser{},
		},
		{
			name: "with groupUsers, ok",
			groupUsers: []*query.GroupUser{
				{
					GroupID:            "group1",
					ResourceOwner:      "org1",
					CreationDate:       timeNow,
					Sequence:           1,
					UserID:             "user1",
					PreferredLoginName: "user1",
					DisplayName:        "user1",
					AvatarUrl:          "example.com/user1.png",
				},
				{
					GroupID:            "group1",
					ResourceOwner:      "org1",
					CreationDate:       timeNow,
					Sequence:           1,
					UserID:             "user2",
					PreferredLoginName: "user2",
					DisplayName:        "user2",
					AvatarUrl:          "example.com/user2.png",
				},
				{
					GroupID:            "group2",
					ResourceOwner:      "org1",
					CreationDate:       timeNow,
					Sequence:           1,
					UserID:             "user1",
					PreferredLoginName: "user1",
					DisplayName:        "user1",
					AvatarUrl:          "example.com/user1.png",
				},
			},
			want: []*group_v2.GroupUser{
				{
					GroupId:        "group1",
					OrganizationId: "org1",
					User: &authorization.User{
						Id:                 "user1",
						DisplayName:        "user1",
						PreferredLoginName: "user1",
						AvatarUrl:          "example.com/user1.png",
						OrganizationId:     "org1",
					},
					CreationDate: timestamppb.New(timeNow),
				},
				{
					GroupId:        "group1",
					OrganizationId: "org1",
					User: &authorization.User{
						Id:                 "user2",
						DisplayName:        "user2",
						PreferredLoginName: "user2",
						AvatarUrl:          "example.com/user2.png",
						OrganizationId:     "org1",
					},
					CreationDate: timestamppb.New(timeNow),
				},
				{
					GroupId:        "group2",
					OrganizationId: "org1",
					User: &authorization.User{
						Id:                 "user1",
						DisplayName:        "user1",
						PreferredLoginName: "user1",
						AvatarUrl:          "example.com/user1.png",
						OrganizationId:     "org1",
					},
					CreationDate: timestamppb.New(timeNow),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := groupUsersToPb(tt.groupUsers)
			assert.Equal(t, tt.want, got)
		})
	}
}

package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	groupUsersStmt = regexp.QuoteMeta(
		"SELECT projections.group_users1.group_id" +
			", projections.group_users1.user_id" +
			", projections.users14_humans.display_name" +
			", projections.login_names3.login_name" +
			", projections.group_users1.resource_owner" +
			", projections.group_users1.instance_id" +
			", projections.users14_humans.avatar_key" +
			", projections.group_users1.creation_date" +
			", projections.group_users1.sequence" +
			", COUNT(*) OVER ()" +
			" FROM projections.group_users1" +
			" LEFT JOIN projections.users14_humans ON projections.group_users1.user_id = projections.users14_humans.user_id AND projections.group_users1.instance_id = projections.users14_humans.instance_id" +
			" LEFT JOIN projections.login_names3 ON projections.group_users1.user_id = projections.login_names3.user_id AND projections.group_users1.instance_id = projections.login_names3.instance_id" +
			" WHERE projections.login_names3.is_primary = $1")

	groupUsersColumns = []string{
		"group_id",
		"user_id",
		"display_name",
		"login_name",
		"resource_owner",
		"instance_id",
		"avatar_key",
		"creation_date",
		"sequence",
		"count",
	}
)

func Test_GroupUsersPrepares(t *testing.T) {
	t.Parallel()
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareGroupUsersQuery no result",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupUsersStmt,
					nil,
					nil,
				),
			},
			object: &GroupUsers{GroupUsers: []*GroupUser{}},
		},
		{
			name:    "prepareGroupUsersQuery with one result",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupUsersStmt,
					groupUsersColumns,
					[][]driver.Value{
						{
							"group-id",
							"user-id",
							"display-name",
							"login-name",
							"resource-owner",
							"instance-id",
							"avatar-key",
							testNow,
							1,
						},
					},
				),
			},
			object: &GroupUsers{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				GroupUsers: []*GroupUser{
					{
						GroupID:            "group-id",
						UserID:             "user-id",
						ResourceOwner:      "resource-owner",
						InstanceID:         "instance-id",
						CreationDate:       testNow,
						Sequence:           1,
						PreferredLoginName: "login-name",
						DisplayName:        "display-name",
						AvatarUrl:          "avatar-key",
					},
				},
			},
		},
		{
			name:    "prepareGroupUsersQuery with multiple results",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupUsersStmt,
					groupUsersColumns,
					[][]driver.Value{
						{
							"group-id-1",
							"user-id-1",
							"display-name-1",
							"login-name-1",
							"resource-owner",
							"instance-id",
							"avatar-key",
							testNow,
							1,
						},
						{
							"group-id-1",
							"user-id-2",
							"display-name-2",
							"login-name-2",
							"resource-owner",
							"instance-id",
							"avatar-key",
							testNow,
							1,
						},
					},
				),
			},
			object: &GroupUsers{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				GroupUsers: []*GroupUser{
					{
						GroupID:            "group-id-1",
						UserID:             "user-id-1",
						ResourceOwner:      "resource-owner",
						InstanceID:         "instance-id",
						CreationDate:       testNow,
						Sequence:           1,
						PreferredLoginName: "login-name-1",
						DisplayName:        "display-name-1",
						AvatarUrl:          "avatar-key",
					},
					{
						GroupID:            "group-id-1",
						UserID:             "user-id-2",
						ResourceOwner:      "resource-owner",
						InstanceID:         "instance-id",
						CreationDate:       testNow,
						Sequence:           1,
						PreferredLoginName: "login-name-2",
						DisplayName:        "display-name-2",
						AvatarUrl:          "avatar-key",
					},
				},
			},
		},
		{
			name:    "prepareGroupUsersQuery sql err",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					groupUsersStmt,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*GroupUsers)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_GroupUsersCheckPermission(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		want        []*GroupUser
		groupUsers  *GroupUsers
		permissions []string
	}{
		{
			name: "no permissions",
			want: []*GroupUser{},
			groupUsers: &GroupUsers{
				GroupUsers: []*GroupUser{
					{GroupID: "group1"}, {GroupID: "group2"}, {GroupID: "group3"},
				},
			},
			permissions: []string{},
		},
		{
			name: "permissions for group1",
			want: []*GroupUser{
				{GroupID: "group1", UserID: "user1"},
			},
			groupUsers: &GroupUsers{
				GroupUsers: []*GroupUser{
					{GroupID: "group1", UserID: "user1"}, {GroupID: "group2", UserID: "user2"}, {GroupID: "group3", UserID: "user3"},
				},
			},
			permissions: []string{"group1"},
		},
		{
			name: "permissions for group2",
			want: []*GroupUser{
				{GroupID: "group2", UserID: "user2"},
			},
			groupUsers: &GroupUsers{
				GroupUsers: []*GroupUser{
					{GroupID: "group1", UserID: "user1"}, {GroupID: "group2", UserID: "user2"}, {GroupID: "group3", UserID: "user3"},
				},
			},
			permissions: []string{"group2"},
		},
		{
			name: "permissions for group1 and group2",
			want: []*GroupUser{
				{GroupID: "group1", UserID: "user1"},
				{GroupID: "group2", UserID: "user1"},
			},
			groupUsers: &GroupUsers{
				GroupUsers: []*GroupUser{
					{GroupID: "group1", UserID: "user1"},
					{GroupID: "group2", UserID: "user1"},
					{GroupID: "group3", UserID: "user1"},
					{GroupID: "group4", UserID: "user2"},
				},
			},
			permissions: []string{"group1", "group2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			checkPermission := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				for _, perm := range tt.permissions {
					if resourceID == perm {
						return nil
					}
				}
				return errors.New("not found")
			}
			groupUsersCheckPermission(context.Background(), tt.groupUsers, checkPermission)
			require.Equal(t, tt.want, tt.groupUsers.GroupUsers)
		})
	}
}

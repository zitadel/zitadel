package query

import (
	"database/sql/driver"
	"regexp"
	"testing"
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
			", projections.group_users1.change_date" +
			", projections.group_users1.sequence" +
			", COUNT(*) OVER ()" +
			" FROM projections.group_users1" +
			" LEFT JOIN projections.users14_humans ON projections.group_users1.user_id = projections.users14_humans.user_id AND projections.group_users1.instance_id = projections.users14_humans.instance_id" +
			" LEFT JOIN projections.login_names3 ON projections.group_users1.user_id = projections.login_names3.user_id AND projections.group_users1.instance_id = projections.login_names3.instance_id" +
			" WHERE projections.login_names3.is_primary = $1")

	groupUsersColumns = []string{
		"group_id",
		"user_id",
		"resource_owner",
		"instance_id",
		"creation_date",
		"change_date",
		"sequence",
		"display_name",
		"avatar_key",
		"login_name",
		"count",
	}
)

func Test_GroupUsersPrepares(t *testing.T) {
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
						ChangeDate:         testNow,
						Sequence:           1,
						PreferredLoginName: "login-name",
						DisplayName:        "display-name",
						AvatarUrl:          "avatar-key",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

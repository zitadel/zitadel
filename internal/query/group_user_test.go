package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
)

var (
	groupUsersQuery = regexp.QuoteMeta("SELECT" +
		"  groupusers.creation_date" +
		", groupusers.change_date" +
		", groupusers.sequence" +
		", groupusers.resource_owner" +
		", groupusers.user_id" +
		", groupusers.group_id" +
		", groupusers.attributes" +
		", projections.login_names3.login_name" +
		", projections.users14_humans.email" +
		", projections.users14_humans.first_name" +
		", projections.users14_humans.last_name" +
		", projections.users14_humans.display_name" +
		", projections.users14_machines.name" +
		", projections.users14_humans.avatar_key" +
		", projections.users14.type" +
		", COUNT(*) OVER () " +
		"FROM projections.group_users AS groupusers " +
		"LEFT JOIN projections.users14_humans " +
		"ON groupusers.user_id = projections.users14_humans.user_id " +
		"AND groupusers.instance_id = projections.users14_humans.instance_id " +
		"LEFT JOIN projections.users14_machines " +
		"ON groupusers.user_id = projections.users14_machines.user_id " +
		"AND groupusers.instance_id = projections.users14_machines.instance_id " +
		"LEFT JOIN projections.users14 " +
		"ON groupusers.user_id = projections.users14.id " +
		"AND groupusers.instance_id = projections.users14.instance_id " +
		"LEFT JOIN projections.login_names3 " +
		"ON groupusers.user_id = projections.login_names3.user_id " +
		"AND groupusers.instance_id = projections.login_names3.instance_id " +
		"WHERE projections.login_names3.is_primary = $1")
	groupUsersColumns = []string{
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"user_id",
		"group_id",
		"roles",
		"login_name",
		"email",
		"first_name",
		"last_name",
		"display_name",
		"name",
		"avatar_key",
		"type",
		"count",
	}
)

func Test_GroupUserPrepares(t *testing.T) {
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
					groupUsersQuery,
					nil,
					nil,
				),
			},
			object: &GroupUsers{
				GroupUsers: []*GroupUser{},
			},
		},
		{
			name:    "prepareGroupUsersQuery human found",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupUsersQuery,
					groupUsersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id",
							"group-id",
							database.TextArray[string]{"role-1", "role-2"},
							"gigi@caos-ag.zitadel.ch",
							"gigi@caos.ch",
							"first-name",
							"last-name",
							"display name",
							nil,
							nil,
							domain.UserTypeHuman,
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
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id",
						GroupID:            "group-id",
						Attributes:         database.TextArray[string]{"role-1", "role-2"},
						PreferredLoginName: "gigi@caos-ag.zitadel.ch",
						Email:              "gigi@caos.ch",
						FirstName:          "first-name",
						LastName:           "last-name",
						DisplayName:        "display name",
						AvatarURL:          "",
						UserType:           domain.UserTypeHuman,
					},
				},
			},
		},
		{
			name:    "prepareGroupUsersQuery machine found",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupUsersQuery,
					groupUsersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id",
							"group-id",
							database.TextArray[string]{"role-1", "role-2"},
							"machine@caos-ag.zitadel.ch",
							nil,
							nil,
							nil,
							nil,
							"machine-name",
							nil,
							domain.UserTypeMachine,
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
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id",
						GroupID:            "group-id",
						Attributes:         database.TextArray[string]{"role-1", "role-2"},
						PreferredLoginName: "machine@caos-ag.zitadel.ch",
						Email:              "",
						FirstName:          "",
						LastName:           "",
						DisplayName:        "machine-name",
						AvatarURL:          "",
						UserType:           domain.UserTypeMachine,
					},
				},
			},
		},
		{
			name:    "prepareGroupUsersQuery multiple users",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupUsersQuery,
					groupUsersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id-1",
							"group-id",
							database.TextArray[string]{"role-1", "role-2"},
							"gigi@caos-ag.zitadel.ch",
							"gigi@caos.ch",
							"first-name",
							"last-name",
							"display name",
							nil,
							nil,
							domain.UserTypeHuman,
						},
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id-2",
							"group-id",
							database.TextArray[string]{"role-1", "role-2"},
							"machine@caos-ag.zitadel.ch",
							nil,
							nil,
							nil,
							nil,
							"machine-name",
							nil,
							domain.UserTypeMachine,
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
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id-1",
						GroupID:            "group-id",
						Attributes:         database.TextArray[string]{"role-1", "role-2"},
						PreferredLoginName: "gigi@caos-ag.zitadel.ch",
						Email:              "gigi@caos.ch",
						FirstName:          "first-name",
						LastName:           "last-name",
						DisplayName:        "display name",
						AvatarURL:          "",
						UserType:           domain.UserTypeHuman,
					},
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id-2",
						GroupID:            "group-id",
						Attributes:         database.TextArray[string]{"role-1", "role-2"},
						PreferredLoginName: "machine@caos-ag.zitadel.ch",
						Email:              "",
						FirstName:          "",
						LastName:           "",
						DisplayName:        "machine-name",
						AvatarURL:          "",
						UserType:           domain.UserTypeMachine,
					},
				},
			},
		},
		{
			name:    "prepareGroupUsersQuery sql err",
			prepare: prepareGroupUsersQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					groupUsersQuery,
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	projectMembersQuery = regexp.QuoteMeta("SELECT" +
		" members.creation_date" +
		", members.change_date" +
		", members.sequence" +
		", members.resource_owner" +
		", members.user_id" +
		", members.roles" +
		", projections.login_names2.login_name" +
		", projections.users8_humans.email" +
		", projections.users8_humans.first_name" +
		", projections.users8_humans.last_name" +
		", projections.users8_humans.display_name" +
		", projections.users8_machines.name" +
		", projections.users8_humans.avatar_key" +
		", COUNT(*) OVER () " +
		"FROM projections.project_members3 AS members " +
		"LEFT JOIN projections.users8_humans " +
		"ON members.user_id = projections.users8_humans.user_id " +
		"AND members.instance_id = projections.users8_humans.instance_id " +
		"LEFT JOIN projections.users8_machines " +
		"ON members.user_id = projections.users8_machines.user_id " +
		"AND members.instance_id = projections.users8_machines.instance_id " +
		"LEFT JOIN projections.login_names2 " +
		"ON members.user_id = projections.login_names2.user_id " +
		"AND members.instance_id = projections.login_names2.instance_id " +
		"WHERE projections.login_names2.is_primary = $1")
	projectMembersColumns = []string{
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"user_id",
		"roles",
		"login_name",
		"email",
		"first_name",
		"last_name",
		"display_name",
		"name",
		"avatar_key",
		"count",
	}
)

func Test_ProjectMemberPrepares(t *testing.T) {
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
			name:    "prepareProjectMembersQuery no result",
			prepare: prepareProjectMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectMembersQuery,
					nil,
					nil,
				),
			},
			object: &Members{
				Members: []*Member{},
			},
		},
		{
			name:    "prepareProjectMembersQuery human found",
			prepare: prepareProjectMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectMembersQuery,
					projectMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id",
							database.StringArray{"role-1", "role-2"},
							"gigi@caos-ag.zitadel.ch",
							"gigi@caos.ch",
							"first-name",
							"last-name",
							"display name",
							nil,
							nil,
						},
					},
				),
			},
			object: &Members{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Members: []*Member{
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id",
						Roles:              database.StringArray{"role-1", "role-2"},
						PreferredLoginName: "gigi@caos-ag.zitadel.ch",
						Email:              "gigi@caos.ch",
						FirstName:          "first-name",
						LastName:           "last-name",
						DisplayName:        "display name",
						AvatarURL:          "",
					},
				},
			},
		},
		{
			name:    "prepareProjectMembersQuery machine found",
			prepare: prepareProjectMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectMembersQuery,
					projectMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id",
							database.StringArray{"role-1", "role-2"},
							"machine@caos-ag.zitadel.ch",
							nil,
							nil,
							nil,
							nil,
							"machine-name",
							nil,
						},
					},
				),
			},
			object: &Members{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Members: []*Member{
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id",
						Roles:              database.StringArray{"role-1", "role-2"},
						PreferredLoginName: "machine@caos-ag.zitadel.ch",
						Email:              "",
						FirstName:          "",
						LastName:           "",
						DisplayName:        "machine-name",
						AvatarURL:          "",
					},
				},
			},
		},
		{
			name:    "prepareProjectMembersQuery multiple users",
			prepare: prepareProjectMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectMembersQuery,
					projectMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id-1",
							database.StringArray{"role-1", "role-2"},
							"gigi@caos-ag.zitadel.ch",
							"gigi@caos.ch",
							"first-name",
							"last-name",
							"display name",
							nil,
							nil,
						},
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"user-id-2",
							database.StringArray{"role-1", "role-2"},
							"machine@caos-ag.zitadel.ch",
							nil,
							nil,
							nil,
							nil,
							"machine-name",
							nil,
						},
					},
				),
			},
			object: &Members{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Members: []*Member{
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id-1",
						Roles:              database.StringArray{"role-1", "role-2"},
						PreferredLoginName: "gigi@caos-ag.zitadel.ch",
						Email:              "gigi@caos.ch",
						FirstName:          "first-name",
						LastName:           "last-name",
						DisplayName:        "display name",
						AvatarURL:          "",
					},
					{
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211206,
						ResourceOwner:      "ro",
						UserID:             "user-id-2",
						Roles:              database.StringArray{"role-1", "role-2"},
						PreferredLoginName: "machine@caos-ag.zitadel.ch",
						Email:              "",
						FirstName:          "",
						LastName:           "",
						DisplayName:        "machine-name",
						AvatarURL:          "",
					},
				},
			},
		},
		{
			name:    "prepareProjectMembersQuery sql err",
			prepare: prepareProjectMembersQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					projectMembersQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

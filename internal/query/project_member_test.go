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
	projectMembersQuery = regexp.QuoteMeta("SELECT" +
		" members.creation_date" +
		", members.change_date" +
		", members.sequence" +
		", members.resource_owner" +
		", members.user_id" +
		", members.roles" +
		", projections.login_names3.login_name" +
		", projections.users9_humans.email" +
		", projections.users9_humans.first_name" +
		", projections.users9_humans.last_name" +
		", projections.users9_humans.display_name" +
		", projections.users9_machines.name" +
		", projections.users9_humans.avatar_key" +
		", projections.users9.type" +
		", COUNT(*) OVER () " +
		"FROM projections.project_members4 AS members " +
		"LEFT JOIN projections.users9_humans " +
		"ON members.user_id = projections.users9_humans.user_id " +
		"AND members.instance_id = projections.users9_humans.instance_id " +
		"LEFT JOIN projections.users9_machines " +
		"ON members.user_id = projections.users9_machines.user_id " +
		"AND members.instance_id = projections.users9_machines.instance_id " +
		"LEFT JOIN projections.users9 " +
		"ON members.user_id = projections.users9.id " +
		"AND members.instance_id = projections.users9.instance_id " +
		"LEFT JOIN projections.login_names3 " +
		"ON members.user_id = projections.login_names3.user_id " +
		"AND members.instance_id = projections.login_names3.instance_id " +
		`AS OF SYSTEM TIME '-1 ms' ` +
		"WHERE projections.login_names3.is_primary = $1")
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
		"type",
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
						Roles:              database.TextArray[string]{"role-1", "role-2"},
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
						Roles:              database.TextArray[string]{"role-1", "role-2"},
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
						Roles:              database.TextArray[string]{"role-1", "role-2"},
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
						Roles:              database.TextArray[string]{"role-1", "role-2"},
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
			object: (*ProjectMembership)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

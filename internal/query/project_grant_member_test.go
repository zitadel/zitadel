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
	projectGrantMembersQuery = regexp.QuoteMeta("SELECT" +
		" members.creation_date" +
		", members.change_date" +
		", members.sequence" +
		", members.resource_owner" +
		", members.user_resource_owner" +
		", members.user_id" +
		", members.roles" +
		", projections.login_names3.login_name" +
		", projections.users14_humans.email" +
		", projections.users14_humans.first_name" +
		", projections.users14_humans.last_name" +
		", projections.users14_humans.display_name" +
		", projections.users14_machines.name" +
		", projections.users14_humans.avatar_key" +
		", projections.users14.type" +
		", COUNT(*) OVER () " +
		"FROM projections.project_grant_members4 AS members " +
		"LEFT JOIN projections.users14_humans " +
		"ON members.user_id = projections.users14_humans.user_id " +
		"AND members.instance_id = projections.users14_humans.instance_id " +
		"LEFT JOIN projections.users14_machines " +
		"ON members.user_id = projections.users14_machines.user_id " +
		"AND members.instance_id = projections.users14_machines.instance_id " +
		"LEFT JOIN projections.users14 " +
		"ON members.user_id = projections.users14.id " +
		"AND members.instance_id = projections.users14.instance_id " +
		"LEFT JOIN projections.login_names3 " +
		"ON members.user_id = projections.login_names3.user_id " +
		"AND members.instance_id = projections.login_names3.instance_id " +
		"LEFT JOIN projections.project_grants4 " +
		"ON members.grant_id = projections.project_grants4.grant_id " +
		"AND members.instance_id = projections.project_grants4.instance_id " +
		`AS OF SYSTEM TIME '-1 ms' ` +
		"WHERE projections.login_names3.is_primary = $1")
	projectGrantMembersColumns = []string{
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"user_resource_owner",
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

func Test_ProjectGrantMemberPrepares(t *testing.T) {
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
			name:    "prepareProjectGrantMembersQuery no result",
			prepare: prepareProjectGrantMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectGrantMembersQuery,
					nil,
					nil,
				),
			},
			object: &Members{
				Members: []*Member{},
			},
		},
		{
			name:    "prepareProjectGrantMembersQuery human found",
			prepare: prepareProjectGrantMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectGrantMembersQuery,
					projectGrantMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"uro",
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
						UserResourceOwner:  "uro",
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
			name:    "prepareProjectGrantMembersQuery machine found",
			prepare: prepareProjectGrantMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectGrantMembersQuery,
					projectGrantMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"uro",
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
						UserResourceOwner:  "uro",
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
			name:    "prepareProjectGrantMembersQuery multiple users",
			prepare: prepareProjectGrantMembersQuery,
			want: want{
				sqlExpectations: mockQueries(
					projectGrantMembersQuery,
					projectGrantMembersColumns,
					[][]driver.Value{
						{
							testNow,
							testNow,
							uint64(20211206),
							"ro",
							"uro",
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
							"uro",
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
						UserResourceOwner:  "uro",
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
						UserResourceOwner:  "uro",
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
			name:    "prepareProjectGrantMembersQuery sql err",
			prepare: prepareProjectGrantMembersQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					projectGrantMembersQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*ProjectGrantMembership)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/lib/pq"
)

var (
	projectGrantMembersQuery = regexp.QuoteMeta("SELECT" +
		" members.creation_date" +
		", members.change_date" +
		", members.sequence" +
		", members.resource_owner" +
		", members.user_id" +
		", members.roles" +
		", zitadel.projections.login_names.login_name" +
		", zitadel.projections.users_humans.email" +
		", zitadel.projections.users_humans.first_name" +
		", zitadel.projections.users_humans.last_name" +
		", zitadel.projections.users_humans.display_name" +
		", zitadel.projections.users_machines.name" +
		", zitadel.projections.users_humans.avatar_key" +
		", COUNT(*) OVER () " +
		"FROM zitadel.projections.project_grant_members as members " +
		"LEFT JOIN zitadel.projections.users_humans " +
		"ON members.user_id = zitadel.projections.users_humans.user_id " +
		"LEFT JOIN zitadel.projections.users_machines " +
		"ON members.user_id = zitadel.projections.users_machines.user_id " +
		"LEFT JOIN zitadel.projections.login_names " +
		"ON members.user_id = zitadel.projections.login_names.user_id " +
		"LEFT JOIN zitadel.projections.project_grants " +
		"ON members.grant_id = zitadel.projections.project_grants.grant_id " +
		"WHERE zitadel.projections.login_names.is_primary = $1")
	projectGrantMembersColumns = []string{
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
							"user-id",
							pq.StringArray{"role-1", "role-2"},
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
						Roles:              []string{"role-1", "role-2"},
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
							"user-id",
							pq.StringArray{"role-1", "role-2"},
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
						Roles:              []string{"role-1", "role-2"},
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
							"user-id-1",
							pq.StringArray{"role-1", "role-2"},
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
							pq.StringArray{"role-1", "role-2"},
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
						Roles:              []string{"role-1", "role-2"},
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
						Roles:              []string{"role-1", "role-2"},
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
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

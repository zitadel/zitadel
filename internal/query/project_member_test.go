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
	projectMembersQuery = regexp.QuoteMeta("SELECT" +
		" zitadel.projections.project_members.creation_date" +
		", zitadel.projections.project_members.change_date" +
		", zitadel.projections.project_members.sequence" +
		", zitadel.projections.project_members.resource_owner" +
		", zitadel.projections.project_members.user_id" +
		", zitadel.projections.project_members.roles" +
		", zitadel.projections.login_names.login_name" +
		", zitadel.projections.users_humans.email" +
		", zitadel.projections.users_humans.first_name" +
		", zitadel.projections.users_humans.last_name" +
		", zitadel.projections.users_humans.display_name" +
		", zitadel.projections.users_machines.name" +
		", zitadel.projections.users_humans.avater_key" +
		", COUNT(*) OVER () " +
		"FROM zitadel.projections.project_members " +
		"LEFT JOIN zitadel.projections.users_humans " +
		"ON zitadel.projections.project_members.user_id = zitadel.projections.users_humans.user_id " +
		"LEFT JOIN zitadel.projections.users_machines " +
		"ON zitadel.projections.project_members.user_id = zitadel.projections.users_machines.user_id " +
		"LEFT JOIN zitadel.projections.login_names " +
		"ON zitadel.projections.project_members.user_id = zitadel.projections.login_names.user_id " +
		"WHERE zitadel.projections.login_names.is_primary = $1")
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
		"avater_key",
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

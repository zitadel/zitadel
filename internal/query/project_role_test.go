package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
)

var (
	prepareProjectRolesStmt = `SELECT projections.project_roles3.project_id,` +
		` projections.project_roles3.creation_date,` +
		` projections.project_roles3.change_date,` +
		` projections.project_roles3.resource_owner,` +
		` projections.project_roles3.sequence,` +
		` projections.project_roles3.role_key,` +
		` projections.project_roles3.display_name,` +
		` projections.project_roles3.group_name,` +
		` COUNT(*) OVER ()` +
		` FROM projections.project_roles3` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareProjectRolesCols = []string{
		"project_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"role_key",
		"display_name",
		"group_name",
		"count",
	}
)

func Test_ProjectRolePrepares(t *testing.T) {
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
			name:    "prepareProjectRolesQuery no result",
			prepare: prepareProjectRolesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareProjectRolesStmt),
					nil,
					nil,
				),
			},
			object: &ProjectRoles{ProjectRoles: []*ProjectRole{}},
		},
		{
			name:    "prepareProjectRolesQuery one result",
			prepare: prepareProjectRolesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareProjectRolesStmt),
					prepareProjectRolesCols,
					[][]driver.Value{
						{
							"project-id",
							testNow,
							testNow,
							"ro",
							uint64(20211111),
							"role-key",
							"role-display-name",
							"role-group",
						},
					},
				),
			},
			object: &ProjectRoles{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				ProjectRoles: []*ProjectRole{
					{
						ProjectID:     "project-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211111,
						Key:           "role-key",
						DisplayName:   "role-display-name",
						Group:         "role-group",
					},
				},
			},
		},
		{
			name:    "prepareProjectRolesQuery multiple result",
			prepare: prepareProjectRolesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareProjectRolesStmt),
					prepareProjectRolesCols,
					[][]driver.Value{
						{
							"project-id",
							testNow,
							testNow,
							"ro",
							uint64(20211111),
							"role-key-1",
							"role-display-name-1",
							"role-group",
						},
						{
							"project-id",
							testNow,
							testNow,
							"ro",
							uint64(20211111),
							"role-key-2",
							"role-display-name-2",
							"role-group",
						},
					},
				),
			},
			object: &ProjectRoles{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				ProjectRoles: []*ProjectRole{
					{
						ProjectID:     "project-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211111,
						Key:           "role-key-1",
						DisplayName:   "role-display-name-1",
						Group:         "role-group",
					},
					{
						ProjectID:     "project-id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211111,
						Key:           "role-key-2",
						DisplayName:   "role-display-name-2",
						Group:         "role-group",
					},
				},
			},
		},
		{
			name:    "prepareProjectRolesQuery sql err",
			prepare: prepareProjectRolesQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareProjectRolesStmt),
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

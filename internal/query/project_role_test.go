package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
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
					regexp.QuoteMeta(`SELECT zitadel.projections.project_roles.project_id,`+
						` zitadel.projections.project_roles.creation_date,`+
						` zitadel.projections.project_roles.change_date,`+
						` zitadel.projections.project_roles.resource_owner,`+
						` zitadel.projections.project_roles.sequence,`+
						` zitadel.projections.project_roles.role_key,`+
						` zitadel.projections.project_roles.display_name,`+
						` zitadel.projections.project_roles.group_name,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.project_roles`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.project_roles.project_id,`+
						` zitadel.projections.project_roles.creation_date,`+
						` zitadel.projections.project_roles.change_date,`+
						` zitadel.projections.project_roles.resource_owner,`+
						` zitadel.projections.project_roles.sequence,`+
						` zitadel.projections.project_roles.role_key,`+
						` zitadel.projections.project_roles.display_name,`+
						` zitadel.projections.project_roles.group_name,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.project_roles`),
					[]string{
						"project_id",
						"creation_date",
						"change_date",
						"resource_owner",
						"sequence",
						"role_key",
						"display_name",
						"group_name",
						"count",
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.project_roles.project_id,`+
						` zitadel.projections.project_roles.creation_date,`+
						` zitadel.projections.project_roles.change_date,`+
						` zitadel.projections.project_roles.resource_owner,`+
						` zitadel.projections.project_roles.sequence,`+
						` zitadel.projections.project_roles.role_key,`+
						` zitadel.projections.project_roles.display_name,`+
						` zitadel.projections.project_roles.group_name,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.project_roles`),
					[]string{
						"project_id",
						"creation_date",
						"change_date",
						"resource_owner",
						"sequence",
						"role_key",
						"display_name",
						"group_name",
						"count",
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.project_roles.project_id,`+
						` zitadel.projections.project_roles.creation_date,`+
						` zitadel.projections.project_roles.change_date,`+
						` zitadel.projections.project_roles.resource_owner,`+
						` zitadel.projections.project_roles.sequence,`+
						` zitadel.projections.project_roles.role_key,`+
						` zitadel.projections.project_roles.display_name,`+
						` zitadel.projections.project_roles.group_name,`+
						` COUNT(*) OVER ()`+
						` FROM zitadel.projections.project_roles`),
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

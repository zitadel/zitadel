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

	"github.com/zitadel/zitadel/internal/domain"
)

var (
	prepareProjectRolesStmt = `SELECT projections.project_roles4.project_id,` +
		` projections.project_roles4.creation_date,` +
		` projections.project_roles4.change_date,` +
		` projections.project_roles4.resource_owner,` +
		` projections.project_roles4.sequence,` +
		` projections.project_roles4.role_key,` +
		` projections.project_roles4.display_name,` +
		` projections.project_roles4.group_name,` +
		` COUNT(*) OVER ()` +
		` FROM projections.project_roles4`
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
			object: (*ProjectRoles)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func Test_projectRolesCheckPermission(t *testing.T) {
	roles := func() []*ProjectRole {
		return []*ProjectRole{
			{ProjectID: "project1", ResourceOwner: "org1", Key: "admin"},
			{ProjectID: "project1", ResourceOwner: "org1", Key: "viewer"},
			{ProjectID: "project2", ResourceOwner: "org1", Key: "admin"},
		}
	}
	tests := []struct {
		name string
		// readableProjects are the project IDs the caller has project.role.read on.
		readableProjects []string
		want             []*ProjectRole
	}{
		{
			name:             "member of one project sees only its roles",
			readableProjects: []string{"project1"},
			want: []*ProjectRole{
				{ProjectID: "project1", ResourceOwner: "org1", Key: "admin"},
				{ProjectID: "project1", ResourceOwner: "org1", Key: "viewer"},
			},
		},
		{
			name:             "member of both projects sees all roles",
			readableProjects: []string{"project1", "project2"},
			want:             roles(),
		},
		{
			name:             "no permission",
			readableProjects: []string{},
			want:             []*ProjectRole{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPermission := func(_ context.Context, permission, _, resourceID string) error {
				if permission != domain.PermissionProjectRoleRead {
					return errors.New("unexpected permission " + permission)
				}
				for _, projectID := range tt.readableProjects {
					if resourceID == projectID {
						return nil
					}
				}
				return errors.New("permission denied")
			}
			projectRoles := &ProjectRoles{ProjectRoles: roles()}
			projectRolesCheckPermission(context.Background(), projectRoles, checkPermission)
			require.Equal(t, tt.want, projectRoles.ProjectRoles)
		})
	}
}

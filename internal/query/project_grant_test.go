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
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	projectGrantsQuery = `SELECT projections.project_grants3.project_id,` +
		` projections.project_grants3.grant_id,` +
		` projections.project_grants3.creation_date,` +
		` projections.project_grants3.change_date,` +
		` projections.project_grants3.resource_owner,` +
		` projections.project_grants3.state,` +
		` projections.project_grants3.sequence,` +
		` projections.projects3.name,` +
		` projections.project_grants3.granted_org_id,` +
		` o.name,` +
		` projections.project_grants3.granted_role_keys,` +
		` r.name,` +
		` COUNT(*) OVER () ` +
		` FROM projections.project_grants3 ` +
		` LEFT JOIN projections.projects3 ON projections.project_grants3.project_id = projections.projects3.id AND projections.project_grants3.instance_id = projections.projects3.instance_id ` +
		` LEFT JOIN projections.orgs AS r ON projections.project_grants3.resource_owner = r.id AND projections.project_grants3.instance_id = r.instance_id` +
		` LEFT JOIN projections.orgs AS o ON projections.project_grants3.granted_org_id = o.id AND projections.project_grants3.instance_id = o.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	projectGrantsCols = []string{
		"project_id",
		"grant_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"name",
		"granted_org_id",
		"name",
		"granted_role_keys",
		"name",
		"count",
	}
	projectGrantQuery = `SELECT projections.project_grants3.project_id,` +
		` projections.project_grants3.grant_id,` +
		` projections.project_grants3.creation_date,` +
		` projections.project_grants3.change_date,` +
		` projections.project_grants3.resource_owner,` +
		` projections.project_grants3.state,` +
		` projections.project_grants3.sequence,` +
		` projections.projects3.name,` +
		` projections.project_grants3.granted_org_id,` +
		` o.name,` +
		` projections.project_grants3.granted_role_keys,` +
		` r.name` +
		` FROM projections.project_grants3 ` +
		` LEFT JOIN projections.projects3 ON projections.project_grants3.project_id = projections.projects3.id AND projections.project_grants3.instance_id = projections.projects3.instance_id ` +
		` LEFT JOIN projections.orgs AS r ON projections.project_grants3.resource_owner = r.id AND projections.project_grants3.instance_id = r.instance_id` +
		` LEFT JOIN projections.orgs AS o ON projections.project_grants3.granted_org_id = o.id AND projections.project_grants3.instance_id = o.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	projectGrantCols = []string{
		"project_id",
		"grant_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"name",
		"granted_org_id",
		"name",
		"granted_role_keys",
		"name",
	}
)

func Test_ProjectGrantPrepares(t *testing.T) {
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
			name:    "prepareProjectGrantsQuery no result",
			prepare: prepareProjectGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectGrantsQuery),
					nil,
					nil,
				),
			},
			object: &ProjectGrants{ProjectGrants: []*ProjectGrant{}},
		},
		{
			name:    "prepareProjectGrantsQuery one result",
			prepare: prepareProjectGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectGrantsQuery),
					projectGrantsCols,
					[][]driver.Value{
						{
							"project-id",
							"grant-id",
							testNow,
							testNow,
							"ro",
							domain.ProjectGrantStateActive,
							20211111,
							"project-name",
							"org-id",
							"org-name",
							database.StringArray{"role-key"},
							"ro-name",
						},
					},
				),
			},
			object: &ProjectGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				ProjectGrants: []*ProjectGrant{
					{
						ProjectID:         "project-id",
						CreationDate:      testNow,
						ChangeDate:        testNow,
						ResourceOwner:     "ro",
						Sequence:          20211111,
						GrantID:           "grant-id",
						State:             domain.ProjectGrantStateActive,
						ProjectName:       "project-name",
						GrantedOrgID:      "org-id",
						OrgName:           "org-name",
						GrantedRoleKeys:   database.StringArray{"role-key"},
						ResourceOwnerName: "ro-name",
					},
				},
			},
		},
		{
			name:    "prepareProjectGrantsQuery no project",
			prepare: prepareProjectGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectGrantsQuery),
					projectGrantsCols,
					[][]driver.Value{
						{
							"project-id",
							"grant-id",
							testNow,
							testNow,
							"ro",
							domain.ProjectGrantStateActive,
							20211111,
							nil,
							"org-id",
							"org-name",
							database.StringArray{"role-key"},
							"ro-name",
						},
					},
				),
			},
			object: &ProjectGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				ProjectGrants: []*ProjectGrant{
					{
						ProjectID:         "project-id",
						CreationDate:      testNow,
						ChangeDate:        testNow,
						ResourceOwner:     "ro",
						Sequence:          20211111,
						GrantID:           "grant-id",
						State:             domain.ProjectGrantStateActive,
						ProjectName:       "",
						GrantedOrgID:      "org-id",
						OrgName:           "org-name",
						GrantedRoleKeys:   database.StringArray{"role-key"},
						ResourceOwnerName: "ro-name",
					},
				},
			},
		},
		{
			name:    "prepareProjectGrantsQuery no org",
			prepare: prepareProjectGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectGrantsQuery),
					projectGrantsCols,
					[][]driver.Value{
						{
							"project-id",
							"grant-id",
							testNow,
							testNow,
							"ro",
							domain.ProjectGrantStateActive,
							20211111,
							"project-name",
							"org-id",
							nil,
							database.StringArray{"role-key"},
							"ro-name",
						},
					},
				),
			},
			object: &ProjectGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				ProjectGrants: []*ProjectGrant{
					{
						ProjectID:         "project-id",
						CreationDate:      testNow,
						ChangeDate:        testNow,
						ResourceOwner:     "ro",
						Sequence:          20211111,
						GrantID:           "grant-id",
						State:             domain.ProjectGrantStateActive,
						ProjectName:       "project-name",
						GrantedOrgID:      "org-id",
						OrgName:           "",
						GrantedRoleKeys:   database.StringArray{"role-key"},
						ResourceOwnerName: "ro-name",
					},
				},
			},
		},
		{
			name:    "prepareProjectGrantsQuery no resource owner",
			prepare: prepareProjectGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectGrantsQuery),
					projectGrantsCols,
					[][]driver.Value{
						{
							"project-id",
							"grant-id",
							testNow,
							testNow,
							"ro",
							domain.ProjectGrantStateActive,
							20211111,
							"project-name",
							"org-id",
							"org-name",
							database.StringArray{"role-key"},
							nil,
						},
					},
				),
			},
			object: &ProjectGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				ProjectGrants: []*ProjectGrant{
					{
						ProjectID:         "project-id",
						CreationDate:      testNow,
						ChangeDate:        testNow,
						ResourceOwner:     "ro",
						Sequence:          20211111,
						GrantID:           "grant-id",
						State:             domain.ProjectGrantStateActive,
						ProjectName:       "project-name",
						GrantedOrgID:      "org-id",
						OrgName:           "org-name",
						GrantedRoleKeys:   database.StringArray{"role-key"},
						ResourceOwnerName: "",
					},
				},
			},
		},
		{
			name:    "prepareProjectGrantsQuery multiple result",
			prepare: prepareProjectGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectGrantsQuery),
					projectGrantsCols,
					[][]driver.Value{
						{
							"project-id",
							"grant-id-1",
							testNow,
							testNow,
							"ro",
							domain.ProjectGrantStateActive,
							20211111,
							"project-name",
							"org-id",
							"org-name",
							database.StringArray{"role-key"},
							"ro-name",
						},
						{
							"project-id",
							"grant-id-2",
							testNow,
							testNow,
							"ro",
							domain.ProjectGrantStateActive,
							20211111,
							"project-name",
							"org-id",
							"org-name",
							database.StringArray{"role-key"},
							"ro-name",
						},
					},
				),
			},
			object: &ProjectGrants{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				ProjectGrants: []*ProjectGrant{
					{
						ProjectID:         "project-id",
						CreationDate:      testNow,
						ChangeDate:        testNow,
						ResourceOwner:     "ro",
						Sequence:          20211111,
						GrantID:           "grant-id-1",
						State:             domain.ProjectGrantStateActive,
						ProjectName:       "project-name",
						GrantedOrgID:      "org-id",
						OrgName:           "org-name",
						GrantedRoleKeys:   database.StringArray{"role-key"},
						ResourceOwnerName: "ro-name",
					},
					{
						ProjectID:         "project-id",
						CreationDate:      testNow,
						ChangeDate:        testNow,
						ResourceOwner:     "ro",
						Sequence:          20211111,
						GrantID:           "grant-id-2",
						State:             domain.ProjectGrantStateActive,
						ProjectName:       "project-name",
						GrantedOrgID:      "org-id",
						OrgName:           "org-name",
						GrantedRoleKeys:   database.StringArray{"role-key"},
						ResourceOwnerName: "ro-name",
					},
				},
			},
		},
		{
			name:    "prepareProjectGrantsQuery sql err",
			prepare: prepareProjectGrantsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(projectGrantsQuery),
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
		{
			name:    "prepareProjectGrantQuery no result",
			prepare: prepareProjectGrantQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectGrantQuery),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*ProjectGrant)(nil),
		},
		{
			name:    "prepareProjectGrantQuery found",
			prepare: prepareProjectGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(projectGrantQuery),
					projectGrantCols,
					[]driver.Value{
						"project-id",
						"grant-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectGrantStateActive,
						20211111,
						"project-name",
						"org-id",
						"org-name",
						database.StringArray{"role-key"},
						"ro-name",
					},
				),
			},
			object: &ProjectGrant{
				ProjectID:         "project-id",
				CreationDate:      testNow,
				ChangeDate:        testNow,
				ResourceOwner:     "ro",
				Sequence:          20211111,
				GrantID:           "grant-id",
				State:             domain.ProjectGrantStateActive,
				ProjectName:       "project-name",
				GrantedOrgID:      "org-id",
				OrgName:           "org-name",
				GrantedRoleKeys:   database.StringArray{"role-key"},
				ResourceOwnerName: "ro-name",
			},
		},
		{
			name:    "prepareProjectGrantQuery no org",
			prepare: prepareProjectGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(projectGrantQuery),
					projectGrantCols,
					[]driver.Value{
						"project-id",
						"grant-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectGrantStateActive,
						20211111,
						"project-name",
						"org-id",
						nil,
						database.StringArray{"role-key"},
						"ro-name",
					},
				),
			},
			object: &ProjectGrant{
				ProjectID:         "project-id",
				CreationDate:      testNow,
				ChangeDate:        testNow,
				ResourceOwner:     "ro",
				Sequence:          20211111,
				GrantID:           "grant-id",
				State:             domain.ProjectGrantStateActive,
				ProjectName:       "project-name",
				GrantedOrgID:      "org-id",
				OrgName:           "",
				GrantedRoleKeys:   database.StringArray{"role-key"},
				ResourceOwnerName: "ro-name",
			},
		},
		{
			name:    "prepareProjectGrantQuery no resource owner",
			prepare: prepareProjectGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(projectGrantQuery),
					projectGrantCols,
					[]driver.Value{
						"project-id",
						"grant-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectGrantStateActive,
						20211111,
						"project-name",
						"org-id",
						"org-name",
						database.StringArray{"role-key"},
						nil,
					},
				),
			},
			object: &ProjectGrant{
				ProjectID:         "project-id",
				CreationDate:      testNow,
				ChangeDate:        testNow,
				ResourceOwner:     "ro",
				Sequence:          20211111,
				GrantID:           "grant-id",
				State:             domain.ProjectGrantStateActive,
				ProjectName:       "project-name",
				GrantedOrgID:      "org-id",
				OrgName:           "org-name",
				GrantedRoleKeys:   database.StringArray{"role-key"},
				ResourceOwnerName: "",
			},
		},
		{
			name:    "prepareProjectGrantQuery no project",
			prepare: prepareProjectGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(projectGrantQuery),
					projectGrantCols,
					[]driver.Value{
						"project-id",
						"grant-id",
						testNow,
						testNow,
						"ro",
						domain.ProjectGrantStateActive,
						20211111,
						nil,
						"org-id",
						"org-name",
						database.StringArray{"role-key"},
						"ro-name",
					},
				),
			},
			object: &ProjectGrant{
				ProjectID:         "project-id",
				CreationDate:      testNow,
				ChangeDate:        testNow,
				ResourceOwner:     "ro",
				Sequence:          20211111,
				GrantID:           "grant-id",
				State:             domain.ProjectGrantStateActive,
				ProjectName:       "",
				GrantedOrgID:      "org-id",
				OrgName:           "org-name",
				GrantedRoleKeys:   database.StringArray{"role-key"},
				ResourceOwnerName: "ro-name",
			},
		},
		{
			name:    "prepareProjectGrantQuery sql err",
			prepare: prepareProjectGrantQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(projectGrantQuery),
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

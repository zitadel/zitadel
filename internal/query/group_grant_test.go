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
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	groupGrantStmt = regexp.QuoteMeta(
		"SELECT projections.group_grants.id" +
			", projections.group_grants.creation_date" +
			", projections.group_grants.change_date" +
			", projections.group_grants.sequence" +
			", projections.group_grants.grant_id" +
			", projections.group_grants.roles" +
			", projections.group_grants.state" +
			", projections.group_grants.group_id" +
			", projections.group.name" +
			", projections.group.description" +
			", projections.group.resource_owner" +
			", projections.group_grants.resource_owner" +
			", projections.orgs1.name" +
			", projections.orgs1.primary_domain" +
			", projections.group_grants.project_id" +
			", projections.projects4.name" +
			", granted_orgs.id" +
			", granted_orgs.name" +
			", granted_orgs.primary_domain" +
			" FROM projections.group_grants" +
			" LEFT JOIN projections.group ON projections.group_grants.group_id = projections.group.id AND projections.group_grants.instance_id = projections.group.instance_id" +
			" LEFT JOIN projections.orgs1 ON projections.group_grants.resource_owner = projections.orgs1.id AND projections.group_grants.instance_id = projections.orgs1.instance_id" +
			" LEFT JOIN projections.projects4 ON projections.group_grants.project_id = projections.projects4.id AND projections.group_grants.instance_id = projections.projects4.instance_id" +
			" LEFT JOIN projections.orgs1 AS granted_orgs ON projections.group.resource_owner = granted_orgs.id AND projections.group.instance_id = granted_orgs.instance_id" +
			` AS OF SYSTEM TIME '-1 ms' `)
	groupGrantCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"grant_id",
		"roles",
		"state",
		"group_id",
		"group_name",
		"group_description",
		"resource_owner", // user resource owner
		"ro",             // user_grant resource owner
		"name",           // org name
		"primary_domain",
		"project_id",
		"name",           // project name
		"id",             // granted org id
		"name",           // granted org name
		"primary_domain", // granted org domain
	}
	groupGrantsStmt = regexp.QuoteMeta(
		"SELECT projections.group_grants.id" +
			", projections.group_grants.creation_date" +
			", projections.group_grants.change_date" +
			", projections.group_grants.sequence" +
			", projections.group_grants.grant_id" +
			", projections.group_grants.roles" +
			", projections.group_grants.state" +
			", projections.group_grants.group_id" +
			", projections.group.name" +
			", projections.group.description" +
			", projections.group.resource_owner" +
			", projections.group_grants.resource_owner" +
			", projections.orgs1.name" +
			", projections.orgs1.primary_domain" +
			", projections.group_grants.project_id" +
			", projections.projects4.name" +
			", granted_orgs.id" +
			", granted_orgs.name" +
			", granted_orgs.primary_domain" +
			", COUNT(*) OVER ()" +
			" FROM projections.group_grants" +
			" LEFT JOIN projections.group ON projections.group_grants.group_id = projections.group.id AND projections.group_grants.instance_id = projections.group.instance_id" +
			" LEFT JOIN projections.orgs1 ON projections.group_grants.resource_owner = projections.orgs1.id AND projections.group_grants.instance_id = projections.orgs1.instance_id" +
			" LEFT JOIN projections.projects4 ON projections.group_grants.project_id = projections.projects4.id AND projections.group_grants.instance_id = projections.projects4.instance_id" +
			" LEFT JOIN projections.orgs1 AS granted_orgs ON projections.group.resource_owner = granted_orgs.id AND projections.group.instance_id = granted_orgs.instance_id" +
			` AS OF SYSTEM TIME '-1 ms' `)
	groupGrantsCols = append(
		groupGrantCols,
		"count",
	)
)

func Test_GroupGrantPrepares(t *testing.T) {
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
			name:    "prepareGroupGrantQuery no result",
			prepare: prepareGroupGrantQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					groupGrantStmt,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*GroupGrant)(nil),
		},
		{
			name:    "prepareGroupGrantQuery found",
			prepare: prepareGroupGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					groupGrantStmt,
					groupGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.GroupGrantStateActive,
						"group-id",
						"groupname",
						"group-description",
						"resource-owner",
						"ro",
						"org-name",
						"primary-domain",
						"project-id",
						"project-name",
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &GroupGrant{
				ID:                 "id",
				CreationDate:       testNow,
				ChangeDate:         testNow,
				Sequence:           20211111,
				GrantID:            "grant-id",
				Roles:              database.TextArray[string]{"role-key"},
				State:              domain.GroupGrantStateActive,
				GroupID:            "group-id",
				GroupName:          "groupname",
				GroupDescription:   "group-description",
				GroupResourceOwner: "resource-owner",
				ResourceOwner:      "ro",
				OrgName:            "org-name",
				OrgPrimaryDomain:   "primary-domain",
				ProjectID:          "project-id",
				ProjectName:        "project-name",
				GrantedOrgID:       "granted-org-id",
				GrantedOrgName:     "granted-org-name",
				GrantedOrgDomain:   "granted-org-domain",
			},
		},
		{
			name:    "prepareGroupGrantQuery (no org) found",
			prepare: prepareGroupGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					groupGrantStmt,
					groupGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.GroupGrantStateActive,
						"group-id",
						"groupname",
						"group-description",
						"resource-owner",
						"ro",
						nil,
						nil,
						"project-id",
						"project-name",
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &GroupGrant{
				ID:                 "id",
				CreationDate:       testNow,
				ChangeDate:         testNow,
				Sequence:           20211111,
				GrantID:            "grant-id",
				Roles:              database.TextArray[string]{"role-key"},
				State:              domain.GroupGrantStateActive,
				GroupID:            "group-id",
				GroupName:          "groupname",
				GroupDescription:   "group-description",
				GroupResourceOwner: "resource-owner",
				ResourceOwner:      "ro",
				OrgName:            "",
				OrgPrimaryDomain:   "",
				ProjectID:          "project-id",
				ProjectName:        "project-name",
				GrantedOrgID:       "granted-org-id",
				GrantedOrgName:     "granted-org-name",
				GrantedOrgDomain:   "granted-org-domain",
			},
		},
		{
			name:    "prepareGroupGrantQuery (no project) found",
			prepare: prepareGroupGrantQuery,
			want: want{
				sqlExpectations: mockQuery(
					groupGrantStmt,
					groupGrantCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						20211111,
						"grant-id",
						database.TextArray[string]{"role-key"},
						domain.GroupGrantStateActive,
						"group-id",
						"groupname",
						"group-description",
						"resource-owner",
						"ro",
						"org-name",
						"primary-domain",
						"project-id",
						nil,
						"granted-org-id",
						"granted-org-name",
						"granted-org-domain",
					},
				),
			},
			object: &GroupGrant{
				ID:                 "id",
				CreationDate:       testNow,
				ChangeDate:         testNow,
				Sequence:           20211111,
				GrantID:            "grant-id",
				Roles:              database.TextArray[string]{"role-key"},
				State:              domain.GroupGrantStateActive,
				GroupID:            "group-id",
				GroupName:          "groupname",
				GroupDescription:   "group-description",
				GroupResourceOwner: "resource-owner",
				ResourceOwner:      "ro",
				OrgName:            "org-name",
				OrgPrimaryDomain:   "primary-domain",
				ProjectID:          "project-id",
				ProjectName:        "",
				GrantedOrgID:       "granted-org-id",
				GrantedOrgName:     "granted-org-name",
				GrantedOrgDomain:   "granted-org-domain",
			},
		},
		{
			name:    "prepareGroupGrantQuery sql err",
			prepare: prepareGroupGrantQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					groupGrantStmt,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*GroupGrant)(nil),
		},
		{
			name:    "prepareGroupGrantsQuery no result",
			prepare: prepareGroupGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupGrantsStmt,
					nil,
					nil,
				),
			},
			object: &GroupGrants{GroupGrants: []*GroupGrant{}},
		},
		{
			name:    "prepareGroupGrantsQuery one grant",
			prepare: prepareGroupGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupGrantsStmt,
					groupGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.GroupGrantStateActive,
							"group-id",
							"groupname",
							"group-description",
							"resource-owner",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &GroupGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				GroupGrants: []*GroupGrant{
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211111,
						GrantID:            "grant-id",
						Roles:              database.TextArray[string]{"role-key"},
						State:              domain.GroupGrantStateActive,
						GroupID:            "group-id",
						GroupName:          "groupname",
						GroupDescription:   "group-description",
						GroupResourceOwner: "resource-owner",
						ResourceOwner:      "ro",
						OrgName:            "org-name",
						OrgPrimaryDomain:   "primary-domain",
						ProjectID:          "project-id",
						ProjectName:        "project-name",
						GrantedOrgID:       "granted-org-id",
						GrantedOrgName:     "granted-org-name",
						GrantedOrgDomain:   "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareGroupGrantsQuery one grant (no org)",
			prepare: prepareGroupGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupGrantsStmt,
					groupGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.GroupGrantStateActive,
							"group-id",
							"groupname",
							"group-description",
							"resource-owner",
							"ro",
							nil,
							nil,
							"project-id",
							"project-name",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &GroupGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				GroupGrants: []*GroupGrant{
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211111,
						GrantID:            "grant-id",
						Roles:              database.TextArray[string]{"role-key"},
						State:              domain.GroupGrantStateActive,
						GroupID:            "group-id",
						GroupName:          "groupname",
						GroupDescription:   "group-description",
						GroupResourceOwner: "resource-owner",
						ResourceOwner:      "ro",
						OrgName:            "",
						OrgPrimaryDomain:   "",
						ProjectID:          "project-id",
						ProjectName:        "project-name",
						GrantedOrgID:       "granted-org-id",
						GrantedOrgName:     "granted-org-name",
						GrantedOrgDomain:   "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareGroupGrantsQuery one grant (no project)",
			prepare: prepareGroupGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupGrantsStmt,
					groupGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.GroupGrantStateActive,
							"group-id",
							"groupname",
							"group-description",
							"resource-owner",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							nil,
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &GroupGrants{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				GroupGrants: []*GroupGrant{
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211111,
						GrantID:            "grant-id",
						Roles:              database.TextArray[string]{"role-key"},
						State:              domain.GroupGrantStateActive,
						GroupID:            "group-id",
						GroupName:          "groupname",
						GroupDescription:   "group-description",
						GroupResourceOwner: "resource-owner",
						ResourceOwner:      "ro",
						OrgName:            "org-name",
						OrgPrimaryDomain:   "primary-domain",
						ProjectID:          "project-id",
						ProjectName:        "",
						GrantedOrgID:       "granted-org-id",
						GrantedOrgName:     "granted-org-name",
						GrantedOrgDomain:   "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareGroupGrantsQuery multiple grants",
			prepare: prepareGroupGrantsQuery,
			want: want{
				sqlExpectations: mockQueries(
					groupGrantsStmt,
					groupGrantsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.GroupGrantStateActive,
							"group-id",
							"groupname",
							"group-description",
							"resource-owner",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
						{
							"id",
							testNow,
							testNow,
							20211111,
							"grant-id",
							database.TextArray[string]{"role-key"},
							domain.GroupGrantStateActive,
							"group-id",
							"groupname",
							"group-description",
							"resource-owner",
							"ro",
							"org-name",
							"primary-domain",
							"project-id",
							"project-name",
							"granted-org-id",
							"granted-org-name",
							"granted-org-domain",
						},
					},
				),
			},
			object: &GroupGrants{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				GroupGrants: []*GroupGrant{
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211111,
						GrantID:            "grant-id",
						Roles:              database.TextArray[string]{"role-key"},
						State:              domain.GroupGrantStateActive,
						GroupID:            "group-id",
						GroupName:          "groupname",
						GroupDescription:   "group-description",
						GroupResourceOwner: "resource-owner",
						ResourceOwner:      "ro",
						OrgName:            "org-name",
						OrgPrimaryDomain:   "primary-domain",
						ProjectID:          "project-id",
						ProjectName:        "project-name",
						GrantedOrgID:       "granted-org-id",
						GrantedOrgName:     "granted-org-name",
						GrantedOrgDomain:   "granted-org-domain",
					},
					{
						ID:                 "id",
						CreationDate:       testNow,
						ChangeDate:         testNow,
						Sequence:           20211111,
						GrantID:            "grant-id",
						Roles:              database.TextArray[string]{"role-key"},
						State:              domain.GroupGrantStateActive,
						GroupID:            "group-id",
						GroupName:          "groupname",
						GroupDescription:   "group-description",
						GroupResourceOwner: "resource-owner",
						ResourceOwner:      "ro",
						OrgName:            "org-name",
						OrgPrimaryDomain:   "primary-domain",
						ProjectID:          "project-id",
						ProjectName:        "project-name",
						GrantedOrgID:       "granted-org-id",
						GrantedOrgName:     "granted-org-name",
						GrantedOrgDomain:   "granted-org-domain",
					},
				},
			},
		},
		{
			name:    "prepareGroupGrantsQuery sql err",
			prepare: prepareGroupGrantsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					groupGrantsStmt,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*GroupGrants)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

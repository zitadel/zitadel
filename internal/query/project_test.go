package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	projectCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"name",
		"project_role_assertion",
		"project_role_check",
		"has_project_check",
		"private_labeling_setting",
	}

	prepareProjectsStmt = `SELECT projections.projects3.id,` +
		` projections.projects3.creation_date,` +
		` projections.projects3.change_date,` +
		` projections.projects3.resource_owner,` +
		` projections.projects3.state,` +
		` projections.projects3.sequence,` +
		` projections.projects3.name,` +
		` projections.projects3.project_role_assertion,` +
		` projections.projects3.project_role_check,` +
		` projections.projects3.has_project_check,` +
		` projections.projects3.private_labeling_setting,` +
		` COUNT(*) OVER ()` +
		` FROM projections.projects3` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareProjectsCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"name",
		"project_role_assertion",
		"project_role_check",
		"has_project_check",
		"private_labeling_setting",
		"count",
	}

	prepareProjectStmt = `SELECT projections.projects3.id,` +
		` projections.projects3.creation_date,` +
		` projections.projects3.change_date,` +
		` projections.projects3.resource_owner,` +
		` projections.projects3.state,` +
		` projections.projects3.sequence,` +
		` projections.projects3.name,` +
		` projections.projects3.project_role_assertion,` +
		` projections.projects3.project_role_check,` +
		` projections.projects3.has_project_check,` +
		` projections.projects3.private_labeling_setting` +
		` FROM projections.projects3` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareProjectCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"state",
		"sequence",
		"name",
		"project_role_assertion",
		"project_role_check",
		"has_project_check",
		"private_labeling_setting",
	}
)

func Test_ProjectPrepares(t *testing.T) {
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
			name:    "prepareProjectsQuery no result",
			prepare: prepareProjectsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareProjectsStmt),
					nil,
					nil,
				),
			},
			object: &Projects{Projects: []*Project{}},
		},
		{
			name:    "prepareProjectsQuery one result",
			prepare: prepareProjectsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareProjectsStmt),
					prepareProjectsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"ro",
							domain.ProjectStateActive,
							uint64(20211108),
							"project-name",
							true,
							true,
							true,
							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
						},
					},
				),
			},
			object: &Projects{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Projects: []*Project{
					{
						ID:                     "id",
						CreationDate:           testNow,
						ChangeDate:             testNow,
						ResourceOwner:          "ro",
						State:                  domain.ProjectStateActive,
						Sequence:               20211108,
						Name:                   "project-name",
						ProjectRoleAssertion:   true,
						ProjectRoleCheck:       true,
						HasProjectCheck:        true,
						PrivateLabelingSetting: domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
					},
				},
			},
		},
		{
			name:    "prepareProjectsQuery multiple result",
			prepare: prepareProjectsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareProjectsStmt),
					prepareProjectsCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							"ro",
							domain.ProjectStateActive,
							uint64(20211108),
							"project-name-1",
							true,
							true,
							true,
							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
						},
						{
							"id-2",
							testNow,
							testNow,
							"ro",
							domain.ProjectStateActive,
							uint64(20211108),
							"project-name-2",
							false,
							false,
							false,
							domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
						},
					},
				),
			},
			object: &Projects{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Projects: []*Project{
					{
						ID:                     "id-1",
						CreationDate:           testNow,
						ChangeDate:             testNow,
						ResourceOwner:          "ro",
						State:                  domain.ProjectStateActive,
						Sequence:               20211108,
						Name:                   "project-name-1",
						ProjectRoleAssertion:   true,
						ProjectRoleCheck:       true,
						HasProjectCheck:        true,
						PrivateLabelingSetting: domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
					},
					{
						ID:                     "id-2",
						CreationDate:           testNow,
						ChangeDate:             testNow,
						ResourceOwner:          "ro",
						State:                  domain.ProjectStateActive,
						Sequence:               20211108,
						Name:                   "project-name-2",
						ProjectRoleAssertion:   false,
						ProjectRoleCheck:       false,
						HasProjectCheck:        false,
						PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					},
				},
			},
		},
		{
			name:    "prepareProjectsQuery sql err",
			prepare: prepareProjectsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareProjectsStmt),
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
			name:    "prepareProjectQuery no result",
			prepare: prepareProjectQuery,
			want: want{
				sqlExpectations: mockQueries(
					prepareProjectStmt,
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
			object: (*Project)(nil),
		},
		{
			name:    "prepareProjectQuery found",
			prepare: prepareProjectQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareProjectStmt),
					prepareProjectCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"ro",
						domain.ProjectStateActive,
						uint64(20211108),
						"project-name",
						true,
						true,
						true,
						domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
					},
				),
			},
			object: &Project{
				ID:                     "id",
				CreationDate:           testNow,
				ChangeDate:             testNow,
				ResourceOwner:          "ro",
				State:                  domain.ProjectStateActive,
				Sequence:               20211108,
				Name:                   "project-name",
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       true,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
			},
		},
		{
			name:    "prepareProjectQuery sql err",
			prepare: prepareProjectQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareProjectStmt),
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

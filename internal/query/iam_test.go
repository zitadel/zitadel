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

func Test_IAMPrepares(t *testing.T) {
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
			name:    "prepareIAMQuery no result",
			prepare: prepareIAMQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.iam.id,`+
						` zitadel.projections.iam.change_date,`+
						` zitadel.projections.iam.sequence,`+
						` zitadel.projections.iam.global_org_id,`+
						` zitadel.projections.iam.iam_project_id,`+
						` zitadel.projections.iam.setup_started,`+
						` zitadel.projections.iam.setup_done`+
						` FROM zitadel.projections.iam`),
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
			object: (*IAM)(nil),
		},
		{
			name:    "prepareIAMQuery found",
			prepare: prepareIAMQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.iam.id,`+
						` zitadel.projections.iam.change_date,`+
						` zitadel.projections.iam.sequence,`+
						` zitadel.projections.iam.global_org_id,`+
						` zitadel.projections.iam.iam_project_id,`+
						` zitadel.projections.iam.setup_started,`+
						` zitadel.projections.iam.setup_done`+
						` FROM zitadel.projections.iam`),
					[]string{
						"id",
						"change_date",
						"sequence",
						"global_org_id",
						"iam_project_id",
						"setup_started",
						"setup_done",
					},
					[]driver.Value{
						"id",
						testNow,
						uint64(20211108),
						"global-org-id",
						"project-id",
						domain.Step2,
						domain.Step1,
					},
				),
			},
			object: &IAM{
				ID:           "id",
				ChangeDate:   testNow,
				Sequence:     20211108,
				GlobalOrgID:  "global-org-id",
				IAMProjectID: "project-id",
				SetupStarted: domain.Step2,
				SetupDone:    domain.Step1,
			},
		},
		{
			name:    "prepareIAMQuery sql err",
			prepare: prepareIAMQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.iam.id,`+
						` zitadel.projections.iam.change_date,`+
						` zitadel.projections.iam.sequence,`+
						` zitadel.projections.iam.global_org_id,`+
						` zitadel.projections.iam.iam_project_id,`+
						` zitadel.projections.iam.setup_started,`+
						` zitadel.projections.iam.setup_done`+
						` FROM zitadel.projections.iam`),
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

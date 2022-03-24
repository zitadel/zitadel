package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
)

func Test_InstancePrepares(t *testing.T) {
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
			name:    "prepareInstanceQuery no result",
			prepare: prepareIAMQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT projections.instance.id,`+
						` projections.instance.change_date,`+
						` projections.instance.sequence,`+
						` projections.instance.global_org_id,`+
						` projections.instance.iam_project_id,`+
						` projections.instance.setup_started,`+
						` projections.instance.setup_done,`+
						` projections.instance.default_language`+
						` FROM projections.instance`),
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
			object: (*Instance)(nil),
		},
		{
			name:    "prepareInstanceQuery found",
			prepare: prepareIAMQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT projections.instance.id,`+
						` projections.instance.change_date,`+
						` projections.instance.sequence,`+
						` projections.instance.global_org_id,`+
						` projections.instance.iam_project_id,`+
						` projections.instance.setup_started,`+
						` projections.instance.setup_done,`+
						` projections.instance.default_language`+
						` FROM projections.instance`),
					[]string{
						"id",
						"change_date",
						"sequence",
						"global_org_id",
						"iam_project_id",
						"setup_started",
						"setup_done",
						"default_language",
					},
					[]driver.Value{
						"id",
						testNow,
						uint64(20211108),
						"global-org-id",
						"project-id",
						domain.Step2,
						domain.Step1,
						"en",
					},
				),
			},
			object: &Instance{
				ID:              "id",
				ChangeDate:      testNow,
				Sequence:        20211108,
				GlobalOrgID:     "global-org-id",
				IAMProjectID:    "project-id",
				SetupStarted:    domain.Step2,
				SetupDone:       domain.Step1,
				DefaultLanguage: language.English,
			},
		},
		{
			name:    "prepareInstanceQuery sql err",
			prepare: prepareIAMQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT projections.instance.id,`+
						` projections.instance.change_date,`+
						` projections.instance.sequence,`+
						` projections.instance.global_org_id,`+
						` projections.instance.iam_project_id,`+
						` projections.instance.setup_started,`+
						` projections.instance.setup_done,`+
						` projections.instance.default_language`+
						` FROM projections.instance`),
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

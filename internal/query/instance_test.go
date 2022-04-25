package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	sq "github.com/Masterminds/squirrel"
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
			name: "prepareInstanceQuery no result",
			prepare: func() (sq.SelectBuilder, func(*sql.Row) (*Instance, error)) {
				return prepareInstanceQuery("")
			},
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT projections.instances.id,`+
						` projections.instances.creation_date,`+
						` projections.instances.change_date,`+
						` projections.instances.sequence,`+
						` projections.instances.global_org_id,`+
						` projections.instances.iam_project_id,`+
						` projections.instances.console_client_id,`+
						` projections.instances.console_app_id,`+
						` projections.instances.setup_started,`+
						` projections.instances.setup_done,`+
						` projections.instances.default_language`+
						` FROM projections.instances`),
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
			name: "prepareInstanceQuery found",
			prepare: func() (sq.SelectBuilder, func(*sql.Row) (*Instance, error)) {
				return prepareInstanceQuery("")
			},
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT projections.instances.id,`+
						` projections.instances.creation_date,`+
						` projections.instances.change_date,`+
						` projections.instances.sequence,`+
						` projections.instances.global_org_id,`+
						` projections.instances.iam_project_id,`+
						` projections.instances.console_client_id,`+
						` projections.instances.console_app_id,`+
						` projections.instances.setup_started,`+
						` projections.instances.setup_done,`+
						` projections.instances.default_language`+
						` FROM projections.instances`),
					[]string{
						"id",
						"creation_date",
						"change_date",
						"sequence",
						"global_org_id",
						"iam_project_id",
						"console_client_id",
						"console_app_id",
						"setup_started",
						"setup_done",
						"default_language",
					},
					[]driver.Value{
						"id",
						testNow,
						testNow,
						uint64(20211108),
						"global-org-id",
						"project-id",
						"client-id",
						"app-id",
						domain.Step2,
						domain.Step1,
						"en",
					},
				),
			},
			object: &Instance{
				ID:           "id",
				CreationDate: testNow,
				ChangeDate:   testNow,
				Sequence:     20211108,
				GlobalOrgID:  "global-org-id",
				IAMProjectID: "project-id",
				ConsoleID:    "client-id",
				ConsoleAppID: "app-id",
				SetupStarted: domain.Step2,
				SetupDone:    domain.Step1,
				DefaultLang:  language.English,
			},
		},
		{
			name: "prepareInstanceQuery sql err",
			prepare: func() (sq.SelectBuilder, func(*sql.Row) (*Instance, error)) {
				return prepareInstanceQuery("")
			},
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT projections.instances.id,`+
						` projections.instances.creation_date,`+
						` projections.instances.change_date,`+
						` projections.instances.sequence,`+
						` projections.instances.global_org_id,`+
						` projections.instances.iam_project_id,`+
						` projections.instances.console_client_id,`+
						` projections.instances.console_app_id,`+
						` projections.instances.setup_started,`+
						` projections.instances.setup_done,`+
						` projections.instances.default_language`+
						` FROM projections.instances`),
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

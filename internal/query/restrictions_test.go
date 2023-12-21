package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/database"
)

var (
	expectedRestrictionsQuery = regexp.QuoteMeta("SELECT projections.restrictions2.aggregate_id," +
		" projections.restrictions2.creation_date," +
		" projections.restrictions2.change_date," +
		" projections.restrictions2.resource_owner," +
		" projections.restrictions2.sequence," +
		" projections.restrictions2.disallow_public_org_registration," +
		" projections.restrictions2.allowed_languages" +
		" FROM projections.restrictions2" +
		" AS OF SYSTEM TIME '-1 ms'",
	)

	restrictionsCols = []string{
		"aggregate_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"disallow_public_org_registration",
		"allowed_languages",
	}
)

func Test_RestrictionsPrepare(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
		object          interface{}
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
	}{
		{
			name:    "prepareRestrictionsQuery no result",
			prepare: prepareRestrictionsQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					expectedRestrictionsQuery,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrNoRows) {
						return fmt.Errorf("err should be sql.ErrNoRows got: %w", err), false
					}
					return nil, true
				},
				object: Restrictions{
					AllowedLanguages: make([]language.Tag, 0),
				},
			},
		},
		{
			name:    "prepareRestrictionsQuery",
			prepare: prepareRestrictionsQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedRestrictionsQuery,
					restrictionsCols,
					[]driver.Value{
						"restrictions1",
						testNow,
						testNow,
						"instance1",
						0,
						true,
						database.TextArray[string]([]string{"en", "de", "ru"}),
					},
				),
				object: Restrictions{
					AggregateID:                   "restrictions1",
					CreationDate:                  testNow,
					ChangeDate:                    testNow,
					ResourceOwner:                 "instance1",
					Sequence:                      0,
					DisallowPublicOrgRegistration: true,
					AllowedLanguages:              []language.Tag{language.Make("en"), language.Make("de"), language.Make("ru")},
				},
			},
		},
		{
			name:    "prepareRestrictionsQuery sql err",
			prepare: prepareRestrictionsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedRestrictionsQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
				object: (*Restrictions)(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.want.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

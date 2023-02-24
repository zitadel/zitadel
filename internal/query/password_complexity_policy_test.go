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
	preparePasswordComplexityPolicyStmt = `SELECT projections.password_complexity_policies2.id,` +
		` projections.password_complexity_policies2.sequence,` +
		` projections.password_complexity_policies2.creation_date,` +
		` projections.password_complexity_policies2.change_date,` +
		` projections.password_complexity_policies2.resource_owner,` +
		` projections.password_complexity_policies2.min_length,` +
		` projections.password_complexity_policies2.has_lowercase,` +
		` projections.password_complexity_policies2.has_uppercase,` +
		` projections.password_complexity_policies2.has_number,` +
		` projections.password_complexity_policies2.has_symbol,` +
		` projections.password_complexity_policies2.is_default,` +
		` projections.password_complexity_policies2.state` +
		` FROM projections.password_complexity_policies2` +
		` AS OF SYSTEM TIME '-1 ms'`
	preparePasswordComplexityPolicyCols = []string{
		"id",
		"sequence",
		"creation_date",
		"change_date",
		"resource_owner",
		"min_length",
		"has_lowercase",
		"has_uppercase",
		"has_number",
		"has_symbol",
		"is_default",
		"state",
	}
)

func Test_PasswordComplexityPolicyPrepares(t *testing.T) {
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
			name:    "preparePasswordComplexityPolicyQuery no result",
			prepare: preparePasswordComplexityPolicyQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(preparePasswordComplexityPolicyStmt),
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
			object: (*PasswordComplexityPolicy)(nil),
		},
		{
			name:    "preparePasswordComplexityPolicyQuery found",
			prepare: preparePasswordComplexityPolicyQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(preparePasswordComplexityPolicyStmt),
					preparePasswordComplexityPolicyCols,
					[]driver.Value{
						"pol-id",
						uint64(20211109),
						testNow,
						testNow,
						"ro",
						8,
						true,
						true,
						true,
						true,
						true,
						domain.PolicyStateActive,
					},
				),
			},
			object: &PasswordComplexityPolicy{
				ID:            "pol-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				Sequence:      20211109,
				ResourceOwner: "ro",
				State:         domain.PolicyStateActive,
				MinLength:     8,
				HasLowercase:  true,
				HasUppercase:  true,
				HasNumber:     true,
				HasSymbol:     true,
				IsDefault:     true,
			},
		},
		{
			name:    "preparePasswordComplexityPolicyQuery sql err",
			prepare: preparePasswordComplexityPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(preparePasswordComplexityPolicyStmt),
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

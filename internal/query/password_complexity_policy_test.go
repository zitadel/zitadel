package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/caos/zitadel/internal/domain"
	errs "github.com/caos/zitadel/internal/errors"
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
					regexp.QuoteMeta(`SELECT zitadel.projections.password_complexity_policies.id,`+
						` zitadel.projections.password_complexity_policies.sequence,`+
						` zitadel.projections.password_complexity_policies.creation_date,`+
						` zitadel.projections.password_complexity_policies.change_date,`+
						` zitadel.projections.password_complexity_policies.resource_owner,`+
						` zitadel.projections.password_complexity_policies.min_length,`+
						` zitadel.projections.password_complexity_policies.has_lowercase,`+
						` zitadel.projections.password_complexity_policies.has_uppercase,`+
						` zitadel.projections.password_complexity_policies.has_number,`+
						` zitadel.projections.password_complexity_policies.has_symbol,`+
						` zitadel.projections.password_complexity_policies.is_default,`+
						` zitadel.projections.password_complexity_policies.state`+
						` FROM zitadel.projections.password_complexity_policies`),
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
					regexp.QuoteMeta(`SELECT zitadel.projections.password_complexity_policies.id,`+
						` zitadel.projections.password_complexity_policies.sequence,`+
						` zitadel.projections.password_complexity_policies.creation_date,`+
						` zitadel.projections.password_complexity_policies.change_date,`+
						` zitadel.projections.password_complexity_policies.resource_owner,`+
						` zitadel.projections.password_complexity_policies.min_length,`+
						` zitadel.projections.password_complexity_policies.has_lowercase,`+
						` zitadel.projections.password_complexity_policies.has_uppercase,`+
						` zitadel.projections.password_complexity_policies.has_number,`+
						` zitadel.projections.password_complexity_policies.has_symbol,`+
						` zitadel.projections.password_complexity_policies.is_default,`+
						` zitadel.projections.password_complexity_policies.state`+
						` FROM zitadel.projections.password_complexity_policies`),
					[]string{
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
					},
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
					regexp.QuoteMeta(`SELECT zitadel.projections.password_complexity_policies.id,`+
						` zitadel.projections.password_complexity_policies.sequence,`+
						` zitadel.projections.password_complexity_policies.creation_date,`+
						` zitadel.projections.password_complexity_policies.change_date,`+
						` zitadel.projections.password_complexity_policies.resource_owner,`+
						` zitadel.projections.password_complexity_policies.min_length,`+
						` zitadel.projections.password_complexity_policies.has_lowercase,`+
						` zitadel.projections.password_complexity_policies.has_uppercase,`+
						` zitadel.projections.password_complexity_policies.has_number,`+
						` zitadel.projections.password_complexity_policies.has_symbol,`+
						` zitadel.projections.password_complexity_policies.is_default,`+
						` zitadel.projections.password_complexity_policies.state`+
						` FROM zitadel.projections.password_complexity_policies`),
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

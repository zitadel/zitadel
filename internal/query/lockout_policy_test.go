package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareLockoutPolicyStmt = `SELECT projections.lockout_policies3.id,` +
		` projections.lockout_policies3.sequence,` +
		` projections.lockout_policies3.creation_date,` +
		` projections.lockout_policies3.change_date,` +
		` projections.lockout_policies3.resource_owner,` +
		` projections.lockout_policies3.show_failure,` +
		` projections.lockout_policies3.max_password_attempts,` +
		` projections.lockout_policies3.max_otp_attempts,` +
		` projections.lockout_policies3.is_default,` +
		` projections.lockout_policies3.state` +
		` FROM projections.lockout_policies3` +
		` AS OF SYSTEM TIME '-1 ms'`

	prepareLockoutPolicyCols = []string{
		"id",
		"sequence",
		"creation_date",
		"change_date",
		"resource_owner",
		"show_failure",
		"max_password_attempts",
		"max_otp_attempts",
		"is_default",
		"state",
	}
)

func Test_LockoutPolicyPrepares(t *testing.T) {
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
			name:    "prepareLockoutPolicyQuery no result",
			prepare: prepareLockoutPolicyQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareLockoutPolicyStmt),
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
			object: (*LockoutPolicy)(nil),
		},
		{
			name:    "prepareLockoutPolicyQuery found",
			prepare: prepareLockoutPolicyQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareLockoutPolicyStmt),
					prepareLockoutPolicyCols,
					[]driver.Value{
						"pol-id",
						uint64(20211109),
						testNow,
						testNow,
						"ro",
						true,
						20,
						20,
						true,
						domain.PolicyStateActive,
					},
				),
			},
			object: &LockoutPolicy{
				ID:                  "pol-id",
				CreationDate:        testNow,
				ChangeDate:          testNow,
				Sequence:            20211109,
				ResourceOwner:       "ro",
				State:               domain.PolicyStateActive,
				ShowFailures:        true,
				MaxPasswordAttempts: 20,
				MaxOTPAttempts:      20,
				IsDefault:           true,
			},
		},
		{
			name:    "prepareLockoutPolicyQuery sql err",
			prepare: prepareLockoutPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareLockoutPolicyStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*LockoutPolicy)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

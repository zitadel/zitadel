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
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT zitadel.projections.lockout_policies.id,`+
						` zitadel.projections.lockout_policies.sequence,`+
						` zitadel.projections.lockout_policies.creation_date,`+
						` zitadel.projections.lockout_policies.change_date,`+
						` zitadel.projections.lockout_policies.resource_owner,`+
						` zitadel.projections.lockout_policies.show_failure,`+
						` zitadel.projections.lockout_policies.max_password_attempts,`+
						` zitadel.projections.lockout_policies.is_default,`+
						` zitadel.projections.lockout_policies.state`+
						` FROM zitadel.projections.lockout_policies`),
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
			object: (*LockoutPolicy)(nil),
		},
		{
			name:    "prepareLockoutPolicyQuery found",
			prepare: prepareLockoutPolicyQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT zitadel.projections.lockout_policies.id,`+
						` zitadel.projections.lockout_policies.sequence,`+
						` zitadel.projections.lockout_policies.creation_date,`+
						` zitadel.projections.lockout_policies.change_date,`+
						` zitadel.projections.lockout_policies.resource_owner,`+
						` zitadel.projections.lockout_policies.show_failure,`+
						` zitadel.projections.lockout_policies.max_password_attempts,`+
						` zitadel.projections.lockout_policies.is_default,`+
						` zitadel.projections.lockout_policies.state`+
						` FROM zitadel.projections.lockout_policies`),
					[]string{
						"id",
						"sequence",
						"creation_date",
						"change_date",
						"resource_owner",
						"show_failure",
						"max_password_attempts",
						"is_default",
						"state",
					},
					[]driver.Value{
						"pol-id",
						uint64(20211109),
						testNow,
						testNow,
						"ro",
						true,
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
				IsDefault:           true,
			},
		},
		{
			name:    "prepareLockoutPolicyQuery sql err",
			prepare: prepareLockoutPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT zitadel.projections.lockout_policies.id,`+
						` zitadel.projections.lockout_policies.sequence,`+
						` zitadel.projections.lockout_policies.creation_date,`+
						` zitadel.projections.lockout_policies.change_date,`+
						` zitadel.projections.lockout_policies.resource_owner,`+
						` zitadel.projections.lockout_policies.show_failure,`+
						` zitadel.projections.lockout_policies.max_password_attempts,`+
						` zitadel.projections.lockout_policies.is_default,`+
						` zitadel.projections.lockout_policies.state`+
						` FROM zitadel.projections.lockout_policies`),
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

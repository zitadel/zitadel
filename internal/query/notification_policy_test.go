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
	notificationPolicyStmt = regexp.QuoteMeta(`SELECT projections.notification_policies.id,` +
		` projections.notification_policies.sequence,` +
		` projections.notification_policies.creation_date,` +
		` projections.notification_policies.change_date,` +
		` projections.notification_policies.resource_owner,` +
		` projections.notification_policies.password_change,` +
		` projections.notification_policies.is_default,` +
		` projections.notification_policies.state` +
		` FROM projections.notification_policies` +
		` AS OF SYSTEM TIME '-1 ms'`)
	notificationPolicyCols = []string{
		"id",
		"sequence",
		"creation_date",
		"change_date",
		"resource_owner",
		"password_change",
		"is_default",
		"state",
	}
)

func Test_NotificationPolicyPrepares(t *testing.T) {
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
			name:    "prepareNotificationPolicyQuery no result",
			prepare: prepareNotificationPolicyQuery,
			want: want{
				sqlExpectations: mockQueries(
					notificationPolicyStmt,
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
						return fmt.Errorf("err should be NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*NotificationPolicy)(nil),
		},
		{
			name:    "prepareNotificationPolicyQuery found",
			prepare: prepareNotificationPolicyQuery,
			want: want{
				sqlExpectations: mockQuery(
					notificationPolicyStmt,
					notificationPolicyCols,
					[]driver.Value{
						"pol-id",
						uint64(20211109),
						testNow,
						testNow,
						"ro",
						true,
						true,
						domain.PolicyStateActive,
					},
				),
			},
			object: &NotificationPolicy{
				ID:             "pol-id",
				CreationDate:   testNow,
				ChangeDate:     testNow,
				Sequence:       20211109,
				ResourceOwner:  "ro",
				State:          domain.PolicyStateActive,
				PasswordChange: true,
				IsDefault:      true,
			},
		},
		{
			name:    "prepareNotificationPolicyQuery sql err",
			prepare: prepareNotificationPolicyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					notificationPolicyStmt,
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

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

func Test_NotificationProviderPrepares(t *testing.T) {
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
			name:    "prepareNotificationProviderQuery no result",
			prepare: prepareDebugNotificationProviderQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT projections.notification_providers2.aggregate_id,`+
						` projections.notification_providers2.creation_date,`+
						` projections.notification_providers2.change_date,`+
						` projections.notification_providers2.sequence,`+
						` projections.notification_providers2.resource_owner,`+
						` projections.notification_providers2.state,`+
						` projections.notification_providers2.provider_type,`+
						` projections.notification_providers2.compact`+
						` FROM projections.notification_providers2`),
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
			object: (*DebugNotificationProvider)(nil),
		},
		{
			name:    "prepareNotificationProviderQuery found",
			prepare: prepareDebugNotificationProviderQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(`SELECT projections.notification_providers2.aggregate_id,`+
						` projections.notification_providers2.creation_date,`+
						` projections.notification_providers2.change_date,`+
						` projections.notification_providers2.sequence,`+
						` projections.notification_providers2.resource_owner,`+
						` projections.notification_providers2.state,`+
						` projections.notification_providers2.provider_type,`+
						` projections.notification_providers2.compact`+
						` FROM projections.notification_providers2`),
					[]string{
						"aggregate_id",
						"creation_date",
						"change_date",
						"sequence",
						"resource_owner",
						"state",
						"provider_type",
						"compact",
					},
					[]driver.Value{
						"agg-id",
						testNow,
						testNow,
						uint64(20211109),
						"ro-id",
						domain.NotificationProviderStateActive,
						domain.NotificationProviderTypeFile,
						true,
					},
				),
			},
			object: &DebugNotificationProvider{
				AggregateID:   "agg-id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				Sequence:      20211109,
				ResourceOwner: "ro-id",
				State:         domain.NotificationProviderStateActive,
				Type:          domain.NotificationProviderTypeFile,
				Compact:       true,
			},
		},
		{
			name:    "prepareNotificationProviderQuery sql err",
			prepare: prepareDebugNotificationProviderQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT projections.notification_providers2.aggregate_id,`+
						` projections.notification_providers2.creation_date,`+
						` projections.notification_providers2.change_date,`+
						` projections.notification_providers2.sequence,`+
						` projections.notification_providers2.resource_owner,`+
						` projections.notification_providers2.state,`+
						` projections.notification_providers2.provider_type,`+
						` projections.notification_providers2.compact`+
						` FROM projections.notification_providers2`),
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

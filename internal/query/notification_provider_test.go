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
	prepareNotificationProviderStmt = `SELECT projections.notification_providers.aggregate_id,` +
		` projections.notification_providers.creation_date,` +
		` projections.notification_providers.change_date,` +
		` projections.notification_providers.sequence,` +
		` projections.notification_providers.resource_owner,` +
		` projections.notification_providers.state,` +
		` projections.notification_providers.provider_type,` +
		` projections.notification_providers.compact` +
		` FROM projections.notification_providers` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareNotificationProviderCols = []string{
		"aggregate_id",
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"state",
		"provider_type",
		"compact",
	}
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
					regexp.QuoteMeta(prepareNotificationProviderStmt),
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
					regexp.QuoteMeta(prepareNotificationProviderStmt),
					prepareNotificationProviderCols,
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
					regexp.QuoteMeta(prepareNotificationProviderStmt),
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

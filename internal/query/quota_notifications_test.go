package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calculateThreshold(t *testing.T) {
	type args struct {
		usedRel             uint16
		notificationPercent uint16
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "80 - below configuration",
			args: args{
				usedRel:             70,
				notificationPercent: 80,
			},
			want: 0,
		},
		{
			name: "80 - below 100 percent use",
			args: args{
				usedRel:             90,
				notificationPercent: 80,
			},
			want: 80,
		},
		{
			name: "80 - above 100 percent use",
			args: args{
				usedRel:             120,
				notificationPercent: 80,
			},
			want: 80,
		},
		{
			name: "80 - more than twice the use",
			args: args{
				usedRel:             190,
				notificationPercent: 80,
			},
			want: 180,
		},
		{
			name: "100 - below 100 percent use",
			args: args{
				usedRel:             90,
				notificationPercent: 100,
			},
			want: 0,
		},
		{
			name: "100 - above 100 percent use",
			args: args{
				usedRel:             120,
				notificationPercent: 100,
			},
			want: 100,
		},
		{
			name: "100 - more than twice the use",
			args: args{
				usedRel:             210,
				notificationPercent: 100,
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateThreshold(tt.args.usedRel, tt.args.notificationPercent)
			assert.Equal(t, int(tt.want), int(got))
		})
	}
}

var (
	expectedQuotaNotificationsQuery = regexp.QuoteMeta(`SELECT projections.quotas_notifications.id,` +
		` projections.quotas_notifications.call_url,` +
		` projections.quotas_notifications.percent,` +
		` projections.quotas_notifications.repeat,` +
		` projections.quotas_notifications.next_due_threshold` +
		` FROM projections.quotas_notifications` +
		` AS OF SYSTEM TIME '-1 ms'`)

	quotaNotificationsCols = []string{
		"id",
		"call_url",
		"percent",
		"repeat",
		"next_due_threshold",
	}
)

func Test_prepareQuotaNotificationsQuery(t *testing.T) {
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
			name:    "prepareQuotaNotificationsQuery no result",
			prepare: prepareQuotaNotificationsQuery,
			want: want{
				sqlExpectations: mockQueries(
					expectedQuotaNotificationsQuery,
					nil,
					nil,
				),
			},
			object: &QuotaNotifications{Configs: []*QuotaNotification{}},
		},
		{
			name:    "prepareQuotaNotificationsQuery",
			prepare: prepareQuotaNotificationsQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedQuotaNotificationsQuery,
					quotaNotificationsCols,
					[]driver.Value{
						"quota-id",
						"url",
						uint16(100),
						true,
						uint16(100),
					},
				),
			},
			object: &QuotaNotifications{
				Configs: []*QuotaNotification{
					{
						ID:               "quota-id",
						CallURL:          "url",
						Percent:          100,
						Repeat:           true,
						NextDueThreshold: 100,
					},
				},
			},
		},
		{
			name:    "prepareQuotaNotificationsQuery sql err",
			prepare: prepareQuotaNotificationsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedQuotaNotificationsQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Quota)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

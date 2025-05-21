package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	expectedQuotaQuery = regexp.QuoteMeta(`SELECT projections.quotas.id,` +
		` projections.quotas.from_anchor,` +
		` projections.quotas.interval,` +
		` projections.quotas.amount,` +
		` projections.quotas.limit_usage,` +
		` now()` +
		` FROM projections.quotas`)

	quotaCols = []string{
		"id",
		"from_anchor",
		"interval",
		"amount",
		"limit_usage",
		"now",
	}
)

func Test_QuotaPrepare(t *testing.T) {
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
			name:    "prepareQuotaQuery no result",
			prepare: prepareQuotaQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					expectedQuotaQuery,
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
			object: (*Quota)(nil),
		},
		{
			name:    "prepareQuotaQuery",
			prepare: prepareQuotaQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedQuotaQuery,
					quotaCols,
					[]driver.Value{
						"quota-id",
						dayNow,
						&pgtype.Interval{
							Days: 1,
						},
						uint64(1000),
						true,
						testNow,
					},
				),
			},
			object: &Quota{
				ID:                 "quota-id",
				From:               dayNow,
				ResetInterval:      time.Hour * 24,
				CurrentPeriodStart: dayNow,
				Amount:             1000,
				Limit:              true,
			},
		},
		{
			name:    "prepareQuotaQuery sql err",
			prepare: prepareQuotaQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedQuotaQuery,
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

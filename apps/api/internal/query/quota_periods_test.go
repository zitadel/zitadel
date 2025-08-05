package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	expectedRemainingQuotaUsageQuery = regexp.QuoteMeta(`SELECT greatest(0, projections.quotas.amount-projections.quotas_periods.usage)` +
		` FROM projections.quotas_periods` +
		` JOIN projections.quotas ON projections.quotas_periods.unit = projections.quotas.unit AND projections.quotas_periods.instance_id = projections.quotas.instance_id`)
	remainingQuotaUsageCols = []string{
		"usage",
	}
)

func Test_prepareRemainingQuotaUsageQuery(t *testing.T) {
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
			name:    "prepareRemainingQuotaUsageQuery no result",
			prepare: prepareRemainingQuotaUsageQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					expectedRemainingQuotaUsageQuery,
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
			object: (*uint64)(nil),
		},
		{
			name:    "prepareRemainingQuotaUsageQuery",
			prepare: prepareRemainingQuotaUsageQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedRemainingQuotaUsageQuery,
					remainingQuotaUsageCols,
					[]driver.Value{
						uint64(100),
					},
				),
			},
			object: uint64P(100),
		},
		{
			name:    "prepareRemainingQuotaUsageQuery sql err",
			prepare: prepareRemainingQuotaUsageQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedRemainingQuotaUsageQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*uint64)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func uint64P(i int) *uint64 {
	u := uint64(i)
	return &u
}

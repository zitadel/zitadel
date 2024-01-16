package query

/*
import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/limits"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	expectedLimitsQuery = regexp.QuoteMeta("SELECT projections.limits2.aggregate_id," +
		" projections.limits2.creation_date," +
		" projections.limits2.change_date," +
		" projections.limits2.resource_owner," +
		" projections.limits2.sequence," +
		" projections.limits2.audit_log_retention," +
		" projections.limits2.block" +
		" FROM projections.limits2" +
		" AS OF SYSTEM TIME '-1 ms'",
	)

	limitsCols = []string{
		"aggregate_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"audit_log_retention",
		"block",
	}
)

func Test_LimitsPrepare(t *testing.T) {
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
			name:    "prepareLimitsQuery no result",
			prepare: prepareLimitsQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					expectedLimitsQuery,
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
			object: (*limits.Limits)(nil),
		},
		{
			name:    "prepareLimitsQuery",
			prepare: prepareLimitsQuery,
			want: want{
				sqlExpectations: mockQuery(
					expectedLimitsQuery,
					limitsCols,
					[]driver.Value{
						"limits1",
						testNow,
						testNow,
						"instance1",
						0,
						intervalDriverValue(t, time.Hour),
						true,
					},
				),
			},
			object: &limits.Limits{
				AggregateID:       "limits1",
				CreationDate:      testNow,
				ChangeDate:        testNow,
				ResourceOwner:     "instance1",
				Sequence:          0,
				AuditLogRetention: gu.Ptr(time.Hour),
				Block:             gu.Ptr(true),
			},
		},
		{
			name:    "prepareLimitsQuery sql err",
			prepare: prepareLimitsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					expectedLimitsQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*limits.Limits)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
*/

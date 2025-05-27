package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
)

var (
	prepareFailedEventsStmt = `SELECT` +
		` projections.failed_events2.projection_name,` +
		` projections.failed_events2.failed_sequence,` +
		` projections.failed_events2.aggregate_type,` +
		` projections.failed_events2.aggregate_id,` +
		` projections.failed_events2.failure_count,` +
		` projections.failed_events2.last_failed,` +
		` projections.failed_events2.error,` +
		` COUNT(*) OVER ()` +
		` FROM projections.failed_events2`

	prepareFailedEventsCols = []string{
		"projection_name",
		"failed_sequence",
		"aggregate_type",
		"aggregate_id",
		"failure_count",
		"last_failed",
		"error",
		"count",
	}
)

func Test_FailedEventsPrepares(t *testing.T) {
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
			name:    "prepareFailedEventsQuery no result",
			prepare: prepareFailedEventsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFailedEventsStmt),
					nil,
					nil,
				),
			},
			object: &FailedEvents{FailedEvents: []*FailedEvent{}},
		},
		{
			name:    "prepareFailedEventsQuery one result",
			prepare: prepareFailedEventsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFailedEventsStmt),
					prepareFailedEventsCols,
					[][]driver.Value{
						{
							"projection-name",
							uint64(20211108),
							"agg-type",
							"agg-id",
							uint64(2),
							testNow,
							"error",
						},
					},
				),
			},
			object: &FailedEvents{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				FailedEvents: []*FailedEvent{
					{
						ProjectionName: "projection-name",
						FailedSequence: 20211108,
						FailureCount:   2,
						LastFailed:     testNow,
						Error:          "error",
						AggregateType:  "agg-type",
						AggregateID:    "agg-id",
					},
				},
			},
		},
		{
			name:    "prepareFailedEventsQuery multiple result",
			prepare: prepareFailedEventsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareFailedEventsStmt),
					prepareFailedEventsCols,
					[][]driver.Value{
						{
							"projection-name",
							uint64(20211108),
							"agg-type",
							"agg-id",
							2,
							testNow,
							"error",
						},
						{
							"projection-name-2",
							uint64(20211108),
							"agg-type",
							"agg-id",
							2,
							nil,
							"error",
						},
					},
				),
			},
			object: &FailedEvents{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				FailedEvents: []*FailedEvent{
					{
						ProjectionName: "projection-name",
						FailedSequence: 20211108,
						FailureCount:   2,
						LastFailed:     testNow,
						Error:          "error",
						AggregateType:  "agg-type",
						AggregateID:    "agg-id",
					},
					{
						ProjectionName: "projection-name-2",
						FailedSequence: 20211108,
						FailureCount:   2,
						Error:          "error",
						AggregateType:  "agg-type",
						AggregateID:    "agg-id",
					},
				},
			},
		},
		{
			name:    "prepareFailedEventsQuery sql err",
			prepare: prepareFailedEventsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareFailedEventsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*FailedEvents)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

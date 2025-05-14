package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/shopspring/decimal"
)

var (
	currentSequenceStmt = `SELECT` +
		` projections.current_states.last_updated,` +
		` projections.current_states.event_date,` +
		` projections.current_states.position,` +
		` projections.current_states.projection_name,` +
		` projections.current_states.aggregate_type,` +
		` projections.current_states.aggregate_id,` +
		` projections.current_states.sequence,` +
		` COUNT(*) OVER ()` +
		` FROM projections.current_states`

	currentSequenceCols = []string{
		"last_updated",
		"event_date",
		"position",
		"projection_name",
		"aggregate_type",
		"aggregate_id",
		"event_sequence",
		"count",
	}
)

func Test_CurrentSequencesPrepares(t *testing.T) {
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
			name:    "prepareCurrentSequencesQuery no result",
			prepare: prepareCurrentStateQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(currentSequenceStmt),
					nil,
					nil,
				),
			},
			object: &CurrentStates{CurrentStates: []*CurrentState{}},
		},
		{
			name:    "prepareCurrentSequencesQuery one result",
			prepare: prepareCurrentStateQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(currentSequenceStmt),
					currentSequenceCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							float64(20211108),
							"projection-name",
							"agg-type",
							"agg-id",
							uint64(20211108),
						},
					},
				),
			},
			object: &CurrentStates{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				CurrentStates: []*CurrentState{
					{
						ProjectionName: "projection-name",
						State: State{
							EventCreatedAt: testNow,
							LastRun:        testNow,
							Position:       decimal.NewFromInt(20211108),
							AggregateID:    "agg-id",
							AggregateType:  "agg-type",
							Sequence:       20211108,
						},
					},
				},
			},
		},
		{
			name:    "prepareCurrentSequencesQuery multiple result",
			prepare: prepareCurrentStateQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(currentSequenceStmt),
					currentSequenceCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							float64(20211108),
							"projection-name",
							"agg-type",
							"agg-id",
							uint64(20211108),
						},
						{
							testNow,
							testNow,
							float64(20211108),
							"projection-name2",
							"agg-type",
							"agg-id",
							uint64(20211108),
						},
					},
				),
			},
			object: &CurrentStates{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				CurrentStates: []*CurrentState{
					{
						ProjectionName: "projection-name",
						State: State{
							EventCreatedAt: testNow,
							Position:       decimal.NewFromInt(20211108),
							LastRun:        testNow,
							AggregateID:    "agg-id",
							AggregateType:  "agg-type",
							Sequence:       20211108,
						},
					},
					{
						ProjectionName: "projection-name2",
						State: State{
							EventCreatedAt: testNow,
							Position:       decimal.NewFromInt(20211108),
							LastRun:        testNow,
							AggregateID:    "agg-id",
							AggregateType:  "agg-type",
							Sequence:       20211108,
						},
					},
				},
			},
		},
		{
			name:    "prepareCurrentSequencesQuery sql err",
			prepare: prepareCurrentStateQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(currentSequenceStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*CurrentStates)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

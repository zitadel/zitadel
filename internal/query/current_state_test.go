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
	currentSequenceStmt = `SELECT` +
		` projections.current_states.last_updated,` +
		` projections.current_states.aggregate_type,` +
		` projections.current_states.aggregate_id,` +
		` projections.current_states.event_date,` +
		` projections.current_states.event_sequence,` +
		` projections.current_states.projection_name,` +
		` COUNT(*) OVER ()` +
		` FROM projections.current_states` +
		" AS OF SYSTEM TIME '-1 ms' "

	currentSequenceCols = []string{
		"last_updated",
		"aggregate_type",
		"aggregate_id",
		"event_date",
		"event_sequence",
		"projection_name",
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
							"agg-type",
							"agg-id",
							testNow,
							uint64(20211108),
							"projection-name",
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
						EventCreationDate: testNow,
						LastRun:           testNow,
						CurrentPosition:   20211108,
						ProjectionName:    "projection-name",
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
							"agg-type",
							"agg-id",
							testNow,
							uint64(20211108),
							"projection-name",
						},
						{
							testNow,
							"agg-type",
							"agg-id",
							testNow,
							uint64(20211108),
							"projection-name2",
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
						EventCreationDate: testNow,
						CurrentPosition:   20211108,
						ProjectionName:    "projection-name",
						LastRun:           testNow,
					},
					{
						EventCreationDate: testNow,
						CurrentPosition:   20211108,
						ProjectionName:    "projection-name2",
						LastRun:           testNow,
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

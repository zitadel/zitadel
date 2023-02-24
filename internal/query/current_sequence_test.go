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
	currentSequenceStmt = `SELECT max(projections.current_sequences.current_sequence) as current_sequence,` +
		` max(projections.current_sequences.timestamp) as timestamp,` +
		` projections.current_sequences.projection_name,` +
		` COUNT(*) OVER ()` +
		` FROM projections.current_sequences` +
		" AS OF SYSTEM TIME '-1 ms' " +
		` GROUP BY projections.current_sequences.projection_name`

	currentSequenceCols = []string{
		"current_sequence",
		"timestamp",
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
			prepare: prepareCurrentSequencesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(currentSequenceStmt),
					nil,
					nil,
				),
			},
			object: &CurrentSequences{CurrentSequences: []*CurrentSequence{}},
		},
		{
			name:    "prepareCurrentSequencesQuery one result",
			prepare: prepareCurrentSequencesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(currentSequenceStmt),
					currentSequenceCols,
					[][]driver.Value{
						{
							uint64(20211108),
							testNow,
							"projection-name",
						},
					},
				),
			},
			object: &CurrentSequences{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				CurrentSequences: []*CurrentSequence{
					{
						Timestamp:       testNow,
						CurrentSequence: 20211108,
						ProjectionName:  "projection-name",
					},
				},
			},
		},
		{
			name:    "prepareCurrentSequencesQuery multiple result",
			prepare: prepareCurrentSequencesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(currentSequenceStmt),
					currentSequenceCols,
					[][]driver.Value{
						{
							uint64(20211108),
							testNow,
							"projection-name",
						},
						{
							uint64(20211108),
							testNow,
							"projection-name-2",
						},
					},
				),
			},
			object: &CurrentSequences{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				CurrentSequences: []*CurrentSequence{
					{
						Timestamp:       testNow,
						CurrentSequence: 20211108,
						ProjectionName:  "projection-name",
					},
					{
						Timestamp:       testNow,
						CurrentSequence: 20211108,
						ProjectionName:  "projection-name-2",
					},
				},
			},
		},
		{
			name:    "prepareCurrentSequencesQuery sql err",
			prepare: prepareCurrentSequencesQuery,
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
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

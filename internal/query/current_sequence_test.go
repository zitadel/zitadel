package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
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
					regexp.QuoteMeta(`SELECT max(projections.current_sequences.current_sequence) as current_sequence,`+
						` max(projections.current_sequences.timestamp) as timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`+
						` GROUP BY projections.current_sequences.projection_name`),
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
					regexp.QuoteMeta(`SELECT max(projections.current_sequences.current_sequence) as current_sequence,`+
						` max(projections.current_sequences.timestamp) as timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`+
						` GROUP BY projections.current_sequences.projection_name`),
					[]string{
						"current_sequence",
						"timestamp",
						"projection_name",
						"count",
					},
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
						ProjectionName:  "project-name",
					},
				},
			},
		},
		{
			name:    "prepareCurrentSequencesQuery multiple result",
			prepare: prepareCurrentSequencesQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT max(projections.current_sequences.current_sequence) as current_sequence,`+
						` max(projections.current_sequences.timestamp) as timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`+
						` GROUP BY projections.current_sequences.projection_name`),
					[]string{
						"current_sequence",
						"timestamp",
						"projection_name",
						"count",
					},
					[][]driver.Value{
						{
							uint64(20211108),
							testNow,
							"project-name",
						},
						{
							uint64(20211108),
							testNow,
							"project-name-2",
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
						ProjectionName:  "project-name",
					},
					{
						Timestamp:       testNow,
						CurrentSequence: 20211108,
						ProjectionName:  "project-name-2",
					},
				},
			},
		},
		{
			name:    "prepareCurrentSequencesQuery sql err",
			prepare: prepareCurrentSequencesQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(`SELECT max(projections.current_sequences.current_sequence) as current_sequence,`+
						` max(projections.current_sequences.timestamp) as timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`+
						` GROUP BY projections.current_sequences.projection_name`),
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

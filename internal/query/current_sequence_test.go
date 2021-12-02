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
					regexp.QuoteMeta(`SELECT projections.current_sequences.aggregate_type,`+
						` projections.current_sequences.current_sequence,`+
						` projections.current_sequences.timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`),
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
					regexp.QuoteMeta(`SELECT projections.current_sequences.aggregate_type,`+
						` projections.current_sequences.current_sequence,`+
						` projections.current_sequences.timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`),
					[]string{
						"aggregate_type",
						"current_sequence",
						"timestamp",
						"projection_name",
						"count",
					},
					[][]driver.Value{
						{
							"type",
							uint64(20211108),
							testNow,
							"project-name",
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
						AggregateType:   "type",
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
					regexp.QuoteMeta(`SELECT projections.current_sequences.aggregate_type,`+
						` projections.current_sequences.current_sequence,`+
						` projections.current_sequences.timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`),
					[]string{
						"aggregate_type",
						"current_sequence",
						"timestamp",
						"projection_name",
						"count",
					},
					[][]driver.Value{
						{
							"type",
							uint64(20211108),
							testNow,
							"project-name",
						},
						{
							"type",
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
						AggregateType:   "type",
						Timestamp:       testNow,
						CurrentSequence: 20211108,
						ProjectionName:  "project-name",
					},
					{
						AggregateType:   "type",
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
					regexp.QuoteMeta(`SELECT projections.current_sequences.aggregate_type,`+
						` projections.current_sequences.current_sequence,`+
						` projections.current_sequences.timestamp,`+
						` projections.current_sequences.projection_name,`+
						` COUNT(*) OVER ()`+
						` FROM projections.current_sequences`),
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

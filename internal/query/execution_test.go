package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareExecutionsStmt = `SELECT projections.executions1.id,` +
		` projections.executions1.change_date,` +
		` projections.executions1.resource_owner,` +
		` projections.executions1.sequence,` +
		` projections.executions1.targets,` +
		` COUNT(*) OVER ()` +
		` FROM projections.executions1`
	prepareExecutionsCols = []string{
		"id",
		"change_date",
		"resource_owner",
		"sequence",
		"targets",
		"count",
	}

	prepareExecutionStmt = `SELECT projections.executions1.id,` +
		` projections.executions1.change_date,` +
		` projections.executions1.resource_owner,` +
		` projections.executions1.sequence,` +
		` projections.executions1.targets` +
		` FROM projections.executions1`
	prepareExecutionCols = []string{
		"id",
		"change_date",
		"resource_owner",
		"sequence",
		"targets",
	}
)

func Test_ExecutionPrepares(t *testing.T) {
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
			name:    "prepareExecutionsQuery no result",
			prepare: prepareExecutionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareExecutionsStmt),
					nil,
					nil,
				),
			},
			object: &Executions{Executions: []*Execution{}},
		},
		{
			name:    "prepareExecutionsQuery one result",
			prepare: prepareExecutionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareExecutionsStmt),
					prepareExecutionsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							"ro",
							uint64(20211109),
							[]byte(`[{"type":1,"target":"include"}]`),
						},
					},
				),
			},
			object: &Executions{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Executions: []*Execution{
					{
						ID: "id",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211109,
						},
						Targets: []*exec.Target{{Type: domain.ExecutionTargetTypeInclude, Target: "include"}},
					},
				},
			},
		},
		{
			name:    "prepareExecutionsQuery multiple result",
			prepare: prepareExecutionsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareExecutionsStmt),
					prepareExecutionsCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							"ro",
							uint64(20211109),
							[]byte(`[{"type":2,"target":"target"}]`),
						},
						{
							"id-2",
							testNow,
							"ro",
							uint64(20211110),
							[]byte(`[{"type":1,"target":"include"}]`),
						},
					},
				),
			},
			object: &Executions{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Executions: []*Execution{
					{
						ID: "id-1",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211109,
						},
						Targets: []*exec.Target{{Type: domain.ExecutionTargetTypeTarget, Target: "target"}},
					},
					{
						ID: "id-2",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211110,
						},
						Targets: []*exec.Target{{Type: domain.ExecutionTargetTypeInclude, Target: "include"}},
					},
				},
			},
		},
		{
			name:    "prepareExecutionsQuery sql err",
			prepare: prepareExecutionsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareExecutionsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Execution)(nil),
		},
		{
			name:    "prepareExecutionQuery no result",
			prepare: prepareExecutionQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareExecutionStmt),
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
			object: (*Execution)(nil),
		},
		{
			name:    "prepareExecutionQuery found",
			prepare: prepareExecutionQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareExecutionStmt),
					prepareExecutionCols,
					[]driver.Value{
						"id",
						testNow,
						"ro",
						uint64(20211109),
						[]byte(`[{"type":2,"target":"target"}]`),
					},
				),
			},
			object: &Execution{
				ID: "id",
				ObjectDetails: domain.ObjectDetails{
					EventDate:     testNow,
					ResourceOwner: "ro",
					Sequence:      20211109,
				},
				Targets: []*exec.Target{{Type: domain.ExecutionTargetTypeTarget, Target: "target"}},
			},
		},
		{
			name:    "prepareExecutionQuery sql err",
			prepare: prepareExecutionQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareExecutionStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Execution)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

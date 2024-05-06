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
	prepareExecutionsStmt = `SELECT projections.executions1.instance_id,` +
		` projections.executions1.id,` +
		` projections.executions1.change_date,` +
		` projections.executions1.sequence,` +
		` execution_targets.targets,` +
		` COUNT(*) OVER ()` +
		` FROM projections.executions1` +
		` JOIN (` +
		`SELECT instance_id, execution_id, JSONB_AGG( JSON_OBJECT( 'position' : position, 'include' : include, 'target' : target_id ) ) as targets` +
		` FROM projections.executions1_targets` +
		` GROUP BY instance_id, execution_id` +
		`)` +
		` AS execution_targets` +
		` ON execution_targets.instance_id = projections.executions1.instance_id` +
		` AND execution_targets.execution_id = projections.executions1.id`
	prepareExecutionsCols = []string{
		"instance_id",
		"id",
		"change_date",
		"sequence",
		"targets",
		"count",
	}

	prepareExecutionStmt = `SELECT projections.executions1.instance_id,` +
		` projections.executions1.id,` +
		` projections.executions1.change_date,` +
		` projections.executions1.sequence,` +
		` execution_targets.targets` +
		` FROM projections.executions1` +
		` JOIN (` +
		`SELECT instance_id, execution_id, JSONB_AGG( JSON_OBJECT( 'position' : position, 'include' : include, 'target' : target_id ) ) as targets` +
		` FROM projections.executions1_targets` +
		` GROUP BY instance_id, execution_id` +
		`)` +
		` AS execution_targets` +
		` ON execution_targets.instance_id = projections.executions1.instance_id` +
		` AND execution_targets.execution_id = projections.executions1.id`
	prepareExecutionCols = []string{
		"instance_id",
		"id",
		"change_date",
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
							"ro",
							"id",
							testNow,
							uint64(20211109),
							[]byte(`[{"position" : 1, "target" : "target"}, {"position" : 2, "include" : "include"}]`),
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
						Targets: []*exec.Target{
							{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
							{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
						},
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
							"ro",
							"id-1",
							testNow,
							uint64(20211109),
							[]byte(`[{"position" : 1, "target" : "target"}, {"position" : 2, "include" : "include"}]`),
						},
						{
							"ro",
							"id-2",
							testNow,
							uint64(20211110),
							[]byte(`[{"position" : 2, "target" : "target"}, {"position" : 1, "include" : "include"}]`),
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
						Targets: []*exec.Target{
							{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
							{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
						},
					},
					{
						ID: "id-2",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211110,
						},
						Targets: []*exec.Target{
							{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
							{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
						},
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
						"ro",
						"id",
						testNow,
						uint64(20211109),
						[]byte(`[{"position" : 1, "target" : "target"}, {"position" : 2, "include" : "include"}]`),
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
				Targets: []*exec.Target{
					{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
				},
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

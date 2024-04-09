package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	executionTable = table{
		name:          projection.ExecutionTable,
		instanceIDCol: projection.ExecutionInstanceIDCol,
	}
	ExecutionColumnID = Column{
		name:  projection.ExecutionIDCol,
		table: executionTable,
	}
	ExecutionColumnCreationDate = Column{
		name:  projection.ExecutionCreationDateCol,
		table: executionTable,
	}
	ExecutionColumnChangeDate = Column{
		name:  projection.ExecutionChangeDateCol,
		table: executionTable,
	}
	ExecutionColumnResourceOwner = Column{
		name:  projection.ExecutionResourceOwnerCol,
		table: executionTable,
	}
	ExecutionColumnInstanceID = Column{
		name:  projection.ExecutionInstanceIDCol,
		table: executionTable,
	}
	ExecutionColumnSequence = Column{
		name:  projection.ExecutionSequenceCol,
		table: executionTable,
	}
	ExecutionColumnTargets = Column{
		name:  projection.ExecutionTargetsCol,
		table: executionTable,
	}
)

var (
	//go:embed execution_targets.sql
	executionTargetsQuery string
	//go:embed execution_targets_combined.sql
	executionTargetsCombinedQuery string
)

type Executions struct {
	SearchResponse
	Executions []*Execution
}

func (e *Executions) SetState(s *State) {
	e.State = s
}

type Execution struct {
	ID string
	domain.ObjectDetails

	Targets []*exec.Target
}

type ExecutionSearchQueries struct {
	SearchRequest
	Queries []SearchQuery
}

func (q *ExecutionSearchQueries) toQuery(query sq.SelectBuilder) sq.SelectBuilder {
	query = q.SearchRequest.toQuery(query)
	for _, q := range q.Queries {
		query = q.toQuery(query)
	}
	return query
}

func (q *Queries) SearchExecutions(ctx context.Context, queries *ExecutionSearchQueries) (executions *Executions, err error) {
	eq := sq.Eq{
		ExecutionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareExecutionsQuery(ctx, q.client)
	return genericRowsQueryWithState[*Executions](ctx, q.client, executionTable, combineToWhereStmt(query, queries.toQuery, eq), scan)
}

func (q *Queries) GetExecutionByID(ctx context.Context, id string) (execution *Execution, err error) {
	eq := sq.Eq{
		ExecutionColumnID.identifier():         id,
		ExecutionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareExecutionQuery(ctx, q.client)
	return genericRowQuery[*Execution](ctx, q.client, query.Where(eq), scan)
}

func NewExecutionInIDsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(ExecutionColumnID, values)
}

func NewExecutionTypeSearchQuery(t domain.ExecutionType) (SearchQuery, error) {
	return NewTextQuery(ExecutionColumnID, t.String(), TextStartsWith)
}

// ExecutionTargets: provide IDs to select all target information,
func (q *Queries) ExecutionTargets(ctx context.Context, ids []string) (execution []*ExecutionTarget, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	if instanceID == "" {
		return nil, nil
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			execution, err = scanExecutionTargets(rows)
			return err
		},
		executionTargetsQuery,
		instanceID,
		strings.Join(ids, ","),
	)
	return execution, err
}

func (q *Queries) ExecutionTargetsCombined(ctx context.Context, ids1, ids2 []string) (execution []*ExecutionTarget, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.End() }()

	instanceID := authz.GetInstance(ctx).InstanceID()
	if instanceID == "" {
		return nil, nil
	}

	err = q.client.QueryContext(ctx,
		func(rows *sql.Rows) error {
			execution, err = scanExecutionTargets(rows)
			return err
		},
		executionTargetsCombinedQuery,
		instanceID,
		strings.Join(ids1, ","),
		strings.Join(ids2, ","),
	)
	return execution, err
}

func prepareExecutionsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*Executions, error)) {
	return sq.Select(
			ExecutionColumnID.identifier(),
			ExecutionColumnChangeDate.identifier(),
			ExecutionColumnResourceOwner.identifier(),
			ExecutionColumnSequence.identifier(),
			ExecutionColumnTargets.identifier(),
			countColumn.identifier(),
		).From(executionTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Executions, error) {
			executions := make([]*Execution, 0)
			var count uint64
			for rows.Next() {
				execution := new(Execution)
				targets := make([]byte, 0)
				err := rows.Scan(
					&execution.ID,
					&execution.EventDate,
					&execution.ResourceOwner,
					&execution.Sequence,
					&targets,
					&count,
				)
				if err != nil {
					return nil, err
				}
				if len(targets) > 0 {
					if err := json.Unmarshal(targets, &execution.Targets); err != nil {
						return nil, err
					}
				}
				executions = append(executions, execution)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-72xfx5jlj7", "Errors.Query.CloseRows")
			}

			return &Executions{
				Executions: executions,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func prepareExecutionQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(row *sql.Row) (*Execution, error)) {
	return sq.Select(
			ExecutionColumnID.identifier(),
			ExecutionColumnChangeDate.identifier(),
			ExecutionColumnResourceOwner.identifier(),
			ExecutionColumnSequence.identifier(),
			ExecutionColumnTargets.identifier(),
		).From(executionTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Execution, error) {
			execution := new(Execution)
			targets := make([]byte, 0)
			err := row.Scan(
				&execution.ID,
				&execution.EventDate,
				&execution.ResourceOwner,
				&execution.Sequence,
				&targets,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-qzn1xycesh", "Errors.Execution.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-f8sjvm4tb8", "Errors.Internal")
			}
			if len(targets) > 0 {
				if err := json.Unmarshal(targets, &execution.Targets); err != nil {
					return nil, err
				}
			}
			return execution, nil
		}
}

type ExecutionTarget struct {
	InstanceID       string
	ExecutionID      string
	TargetID         string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
}

func (e *ExecutionTarget) GetExecutionID() string {
	return e.ExecutionID
}
func (e *ExecutionTarget) GetTargetID() string {
	return e.TargetID
}
func (e *ExecutionTarget) IsInterruptOnError() bool {
	return e.InterruptOnError
}
func (e *ExecutionTarget) GetEndpoint() string {
	return e.Endpoint
}
func (e *ExecutionTarget) GetTargetType() domain.TargetType {
	return e.TargetType
}
func (e *ExecutionTarget) GetTimeout() time.Duration {
	return e.Timeout
}

func scanExecutionTargets(rows *sql.Rows) ([]*ExecutionTarget, error) {
	targets := make([]*ExecutionTarget, 0)
	for rows.Next() {
		target := new(ExecutionTarget)

		var (
			instanceID       = &sql.NullString{}
			executionID      = &sql.NullString{}
			targetID         = &sql.NullString{}
			targetType       = &sql.NullInt32{}
			endpoint         = &sql.NullString{}
			timeout          = &sql.NullInt64{}
			interruptOnError = &sql.NullBool{}
		)

		err := rows.Scan(
			instanceID,
			executionID,
			targetID,
			targetType,
			endpoint,
			timeout,
			interruptOnError,
		)

		if err != nil {
			return nil, err
		}

		target.InstanceID = instanceID.String
		target.ExecutionID = executionID.String
		target.TargetID = targetID.String
		target.TargetType = domain.TargetType(targetType.Int32)
		target.Endpoint = endpoint.String
		target.Timeout = time.Duration(timeout.Int64)
		target.InterruptOnError = interruptOnError.Bool

		targets = append(targets, target)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-37ardr0pki", "Errors.Query.CloseRows")
	}

	return targets, nil
}

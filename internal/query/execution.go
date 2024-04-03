package query

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
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
	ExecutionColumnIncludes = Column{
		name:  projection.ExecutionIncludesCol,
		table: executionTable,
	}
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

	Targets  database.TextArray[string]
	Includes database.TextArray[string]
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

func NewExecutionTargetSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ExecutionColumnTargets, value, TextListContains)
}

func NewExecutionIncludeSearchQuery(value string) (SearchQuery, error) {
	return NewTextQuery(ExecutionColumnIncludes, value, TextListContains)
}

func prepareExecutionsQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(rows *sql.Rows) (*Executions, error)) {
	return sq.Select(
			ExecutionColumnID.identifier(),
			ExecutionColumnChangeDate.identifier(),
			ExecutionColumnResourceOwner.identifier(),
			ExecutionColumnSequence.identifier(),
			ExecutionColumnTargets.identifier(),
			ExecutionColumnIncludes.identifier(),
			countColumn.identifier(),
		).From(executionTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*Executions, error) {
			executions := make([]*Execution, 0)
			var count uint64
			for rows.Next() {
				execution := new(Execution)
				err := rows.Scan(
					&execution.ID,
					&execution.EventDate,
					&execution.ResourceOwner,
					&execution.Sequence,
					&execution.Targets,
					&execution.Includes,
					&count,
				)
				if err != nil {
					return nil, err
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
			ExecutionColumnIncludes.identifier(),
		).From(executionTable.identifier()).
			PlaceholderFormat(sq.Dollar),
		func(row *sql.Row) (*Execution, error) {
			execution := new(Execution)
			err := row.Scan(
				&execution.ID,
				&execution.EventDate,
				&execution.ResourceOwner,
				&execution.Sequence,
				&execution.Targets,
				&execution.Includes,
			)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-qzn1xycesh", "Errors.Execution.NotFound")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-f8sjvm4tb8", "Errors.Internal")
			}
			return execution, nil
		}
}

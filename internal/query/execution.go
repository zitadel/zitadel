package query

import (
	"cmp"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"slices"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	exec "github.com/zitadel/zitadel/internal/repository/execution"
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
	ExecutionColumnInstanceID = Column{
		name:  projection.ExecutionInstanceIDCol,
		table: executionTable,
	}
	executionTargetsTable = table{
		name:          projection.ExecutionTable + "_" + projection.ExecutionTargetSuffix,
		instanceIDCol: projection.ExecutionTargetInstanceIDCol,
	}
	executionTargetsTableAlias       = executionTargetsTable.setAlias("execution_targets")
	ExecutionTargetsColumnInstanceID = Column{
		name:  projection.ExecutionTargetInstanceIDCol,
		table: executionTargetsTableAlias,
	}
	ExecutionTargetsColumnExecutionID = Column{
		name:  projection.ExecutionTargetExecutionIDCol,
		table: executionTargetsTableAlias,
	}
	executionTargetsListCol = Column{
		name:  "targets",
		table: executionTargetsTableAlias,
	}
)

var (
	//go:embed execution_targets.sql
	executionTargetsQuery string
)

type Executions struct {
	SearchResponse
	Executions []*Execution
}

func (e *Executions) SetState(s *State) {
	e.State = s
}

type Execution struct {
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
	query, scan := prepareExecutionsQuery()
	return genericRowsQueryWithState(ctx, q.client, executionTable, combineToWhereStmt(query, queries.toQuery, eq), scan)
}

func (q *Queries) GetExecutionByID(ctx context.Context, id string) (execution *Execution, err error) {
	eq := sq.Eq{
		ExecutionColumnID.identifier():         id,
		ExecutionColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, scan := prepareExecutionQuery()
	return genericRowQuery(ctx, q.client, query.Where(eq), scan)
}

func NewExecutionInIDsSearchQuery(values []string) (SearchQuery, error) {
	return NewInTextQuery(ExecutionColumnID, values)
}

func NewExecutionTypeSearchQuery(t domain.ExecutionType) (SearchQuery, error) {
	return NewTextQuery(ExecutionColumnID, t.String(), TextStartsWith)
}

func NewTargetSearchQuery(target string) (SearchQuery, error) {
	data, err := targetItemJSONB(domain.ExecutionTargetTypeTarget, target)
	if err != nil {
		return nil, err
	}
	return NewListContains(executionTargetsListCol, data)
}

func NewIncludeSearchQuery(include string) (SearchQuery, error) {
	data, err := targetItemJSONB(domain.ExecutionTargetTypeInclude, include)
	if err != nil {
		return nil, err
	}
	return NewListContains(executionTargetsListCol, data)
}

// marshall executionTargets into the same JSONB structure as in the SQL queries
func targetItemJSONB(t domain.ExecutionTargetType, targetItem string) ([]byte, error) {
	var target *executionTarget
	switch t {
	case domain.ExecutionTargetTypeTarget:
		target = &executionTarget{Target: targetItem}
	case domain.ExecutionTargetTypeInclude:
		target = &executionTarget{Include: targetItem}
	case domain.ExecutionTargetTypeUnspecified:
		return nil, nil
	default:
		return nil, nil
	}
	return json.Marshal([]*executionTarget{target})
}

func prepareExecutionQuery() (sq.SelectBuilder, func(row *sql.Row) (*Execution, error)) {
	return sq.Select(
			ExecutionColumnInstanceID.identifier(),
			ExecutionColumnID.identifier(),
			ExecutionColumnCreationDate.identifier(),
			ExecutionColumnChangeDate.identifier(),
			executionTargetsListCol.identifier(),
		).From(executionTable.identifier()).
			Join("(" + executionTargetsQuery + ") AS " + executionTargetsTableAlias.alias + " ON " +
				ExecutionTargetsColumnInstanceID.identifier() + " = " + ExecutionColumnInstanceID.identifier() + " AND " +
				ExecutionTargetsColumnExecutionID.identifier() + " = " + ExecutionColumnID.identifier(),
			).
			PlaceholderFormat(sq.Dollar),
		scanExecution
}

func prepareExecutionsQuery() (sq.SelectBuilder, func(rows *sql.Rows) (*Executions, error)) {
	return sq.Select(
			ExecutionColumnInstanceID.identifier(),
			ExecutionColumnID.identifier(),
			ExecutionColumnCreationDate.identifier(),
			ExecutionColumnChangeDate.identifier(),
			executionTargetsListCol.identifier(),
			countColumn.identifier(),
		).From(executionTable.identifier()).
			Join("(" + executionTargetsQuery + ") AS " + executionTargetsTableAlias.alias + " ON " +
				ExecutionTargetsColumnInstanceID.identifier() + " = " + ExecutionColumnInstanceID.identifier() + " AND " +
				ExecutionTargetsColumnExecutionID.identifier() + " = " + ExecutionColumnID.identifier(),
			).
			PlaceholderFormat(sq.Dollar),
		scanExecutions
}

type executionTarget struct {
	Position int    `json:"position,omitempty"`
	Include  string `json:"include,omitempty"`
	Target   string `json:"target,omitempty"`
}

func scanExecution(row *sql.Row) (*Execution, error) {
	execution := new(Execution)
	targets := make([]byte, 0)

	err := row.Scan(
		&execution.ResourceOwner,
		&execution.ID,
		&execution.CreationDate,
		&execution.EventDate,
		&targets,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowNotFound(err, "QUERY-qzn1xycesh", "Errors.Execution.NotFound")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-f8sjvm4tb8", "Errors.Internal")
	}

	executionTargets := make([]*executionTarget, 0)
	if err := json.Unmarshal(targets, &executionTargets); err != nil {
		return nil, err
	}

	execution.Targets = make([]*exec.Target, len(executionTargets))
	for i := range executionTargets {
		if executionTargets[i].Target != "" {
			execution.Targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeTarget, Target: executionTargets[i].Target}
		}
		if executionTargets[i].Include != "" {
			execution.Targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeInclude, Target: executionTargets[i].Include}
		}
	}

	return execution, nil
}

func executionTargetsUnmarshal(data []byte) ([]*exec.Target, error) {
	executionTargets := make([]*executionTarget, 0)
	if err := json.Unmarshal(data, &executionTargets); err != nil {
		return nil, err
	}

	targets := make([]*exec.Target, len(executionTargets))
	slices.SortFunc(executionTargets, func(a, b *executionTarget) int {
		return cmp.Compare(a.Position, b.Position)
	})
	for i, item := range executionTargets {
		if item.Target != "" {
			targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeTarget, Target: item.Target}
		}
		if item.Include != "" {
			targets[i] = &exec.Target{Type: domain.ExecutionTargetTypeInclude, Target: item.Include}
		}
	}
	return targets, nil
}

func scanExecutions(rows *sql.Rows) (*Executions, error) {
	executions := make([]*Execution, 0)
	var count uint64

	for rows.Next() {
		execution := new(Execution)
		targets := make([]byte, 0)

		err := rows.Scan(
			&execution.ResourceOwner,
			&execution.ID,
			&execution.CreationDate,
			&execution.EventDate,
			&targets,
			&count,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, zerrors.ThrowNotFound(err, "QUERY-tbrmno85vp", "Errors.Execution.NotFound")
			}
			return nil, zerrors.ThrowInternal(err, "QUERY-tyw2ydsj84", "Errors.Internal")
		}

		execution.Targets, err = executionTargetsUnmarshal(targets)
		if err != nil {
			return nil, err
		}
		executions = append(executions, execution)
	}

	if err := rows.Close(); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-yhka3fs3mw", "Errors.Query.CloseRows")
	}

	return &Executions{
		Executions: executions,
		SearchResponse: SearchResponse{
			Count: count,
		},
	}, nil
}
